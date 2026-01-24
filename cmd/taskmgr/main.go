package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/hiroto-aibara/secretary-ai/internal/handler"
	"github.com/hiroto-aibara/secretary-ai/internal/infra/watcher"
	yamlstore "github.com/hiroto-aibara/secretary-ai/internal/infra/yaml"
	"github.com/hiroto-aibara/secretary-ai/internal/usecase"
	"github.com/hiroto-aibara/secretary-ai/web"
)

func main() {
	basePath := ".tasks"

	// infra
	store := yamlstore.NewStore(basePath)
	cardRepo := yamlstore.NewCardRepositoryAdapter(store)
	hub := handler.NewHub()
	w := watcher.New(hub, basePath)

	// usecase
	boardUC := usecase.NewBoardUseCase(store)
	cardUC := usecase.NewCardUseCase(cardRepo, store)

	// handler
	boardH := handler.NewBoardHandler(boardUC)
	cardH := handler.NewCardHandler(cardUC)
	wsH := handler.NewWSHandler(hub)

	// router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	boardH.Register(r)
	cardH.Register(r)
	wsH.Register(r)

	// static files (embedded frontend)
	r.Get("/*", web.SPAHandler())

	// start watcher
	watchCtx, watchCancel := context.WithCancel(context.Background())
	defer watchCancel()
	go func() {
		if err := w.Start(watchCtx); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("watcher failed", "error", err)
		}
	}()

	// server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		slog.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	watchCancel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}
}
