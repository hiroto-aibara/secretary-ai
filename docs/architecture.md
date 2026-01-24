# アーキテクチャ

## 概要

YAML駆動のカンバン式タスク管理ツール。
タスクの実態はYAMLファイルで、Claude Code（CLI）による直接編集と、Web UIからのREST API経由操作の両方に対応する。

## システム構成

```
┌─────────────────┐            ┌──────────────────┐
│  Claude Code    │── edit ──▶ │  .tasks/*.yaml   │ ◀── source of truth
└─────────────────┘            └────────┬─────────┘
                                        │ fsnotify (watch)
┌─────────────────┐            ┌────────▼─────────┐
│  Web UI         │◀── WS ───▶│  Go Server       │
│  (React Kanban) │◀── REST ──▶│  (chi + embed)   │
└─────────────────┘            └──────────────────┘
```

- **Claude Code**: YAMLファイルを直接編集。ファイル監視によりサーバーが変更を自動検知
- **Web UI**: REST API経由でカードCRUD操作。WebSocketでリアルタイム更新を受信
- **Go Server**: YAMLの読み書き、REST API提供、フロントエンド静的配信

## データ構造

### ディレクトリレイアウト

```
.tasks/
├── config.yaml              # グローバル設定
└── boards/
    ├── project-alpha/
    │   ├── board.yaml       # ボードメタ（名前・リスト定義・順序）
    │   └── cards/
    │       ├── 20260124-001.yaml
    │       └── 20260124-002.yaml
    └── project-beta/
        ├── board.yaml
        └── cards/
            └── ...
```

### スキーマ

#### config.yaml

```yaml
default_board: project-alpha
```

#### board.yaml

```yaml
id: project-alpha
name: "Project Alpha"
lists:
  - id: todo
    name: "Todo"
  - id: in-progress
    name: "In Progress"
  - id: done
    name: "Done"
```

#### カードYAML（例: 20260124-001.yaml）

```yaml
id: "20260124-001"
title: "ログイン機能の実装"
list: in-progress
order: 0
description: |
  OAuth2を使った認証フロー
labels:
  - feature
  - auth
archived: false
created_at: 2026-01-24T10:00:00+09:00
updated_at: 2026-01-24T15:00:00+09:00
```

## アーカイブ仕様

- カードYAMLの `archived: true` フラグで管理
- アーカイブされたカードはカンバン上のリストに表示されない
- Web UIにアーカイブ一覧ビューを用意（フィルタ切り替え）
- 復元操作で `archived: false` に戻し、元のリストに復帰

## バックエンドレイヤー設計

### レイヤー構成と依存方向

```
handler → usecase → domain ← infra
```

- `domain`: エンティティ + リポジトリインターフェース。他のどのパッケージにも依存しない
- `usecase`: ビジネスロジック。`domain` のインターフェースにのみ依存
- `handler`: HTTPリクエスト/レスポンス処理。`usecase` に依存
- `infra`: `domain` インターフェースの具体実装（YAML永続化、ファイル監視等）

### ルール

| # | ルール | 詳細 |
|---|--------|------|
| 1 | 依存方向は内側のみ | handler→usecase→domain の方向のみ許可。逆方向の import は禁止 |
| 2 | インターフェース定義は domain | リポジトリ等のインターフェースは `domain` パッケージに定義 |
| 3 | DI はコンストラクタ注入 | フレームワーク不使用。`New*` 関数で依存を受け取る |
| 4 | DI 配線は main.go に集約 | `cmd/taskmgr/main.go` が唯一の配線ポイント |
| 5 | handler はロジックを持たない | リクエスト解析 + usecase 呼び出し + レスポンス構築のみ |
| 6 | usecase は infra を知らない | インターフェース経由でのみデータアクセス |
| 7 | Notifier は infra 内で完結 | ファイル監視による通知は infra/watcher 内で自己完結。UseCase に注入しない |

### Notifier（WebSocket通知）の設計方針

UseCaseにNotifierを注入**しない**。理由:

- CLI（YAML直接編集）と API（UseCase経由）の両方がYAMLファイル変更に収束するため、fsnotify による一元的な検知で通知が完結する
- UseCase に通知責務を持たせると、二重通知やCLI経路との不整合が生じる

```
パス1: Claude Code → YAML編集 → fsnotify → WebSocket通知
パス2: Web UI → handler → usecase → repository(YAML書込) → fsnotify → WebSocket通知
```

### インターフェース定義例

```go
// domain/repository.go
package domain

import "context"

type BoardRepository interface {
    List(ctx context.Context) ([]Board, error)
    Get(ctx context.Context, id string) (*Board, error)
    Save(ctx context.Context, board *Board) error
    Delete(ctx context.Context, id string) error
}

type CardRepository interface {
    ListByBoard(ctx context.Context, boardID string, includeArchived bool) ([]Card, error)
    Get(ctx context.Context, boardID, cardID string) (*Card, error)
    Save(ctx context.Context, boardID string, card *Card) error
    Delete(ctx context.Context, boardID, cardID string) error
}
```

### DI配線例（main.go）

```go
func main() {
    // infra
    yamlStore := yaml.NewStore(".tasks")
    wsHub := websocket.NewHub()
    watcher := watcher.New(wsHub, ".tasks")

    // usecase
    boardUC := usecase.NewBoardUseCase(yamlStore)
    cardUC := usecase.NewCardUseCase(yamlStore, yamlStore)

    // handler
    boardH := handler.NewBoardHandler(boardUC)
    cardH := handler.NewCardHandler(cardUC)
    wsH := handler.NewWSHandler(wsHub)

    // router
    r := chi.NewRouter()
    boardH.Register(r)
    cardH.Register(r)
    wsH.Register(r)

    go watcher.Start()
    http.ListenAndServe(":8080", r)
}
```

## 技術スタック

### バックエンド

| ライブラリ | 用途 |
|-----------|------|
| Go 1.22+  | 言語 |
| github.com/go-chi/chi/v5 | HTTPルーター |
| gopkg.in/yaml.v3 | YAML読み書き |
| github.com/fsnotify/fsnotify | ファイル監視 |
| github.com/gorilla/websocket | WebSocket |

### フロントエンド

| ライブラリ | 用途 |
|-----------|------|
| React 18  | UI |
| TypeScript | 型安全性 |
| Vite      | ビルドツール |
| @dnd-kit/core + sortable | ドラッグ&ドロップ |
| CSS Modules | スタイリング |

## ソースコード構造

```
SecretaryAi/
├── cmd/
│   └── taskmgr/
│       └── main.go           # エントリポイント（DI配線 + サーバー起動）
├── internal/
│   ├── domain/
│   │   ├── board.go          # Board, List エンティティ
│   │   ├── card.go           # Card エンティティ
│   │   └── repository.go    # BoardRepository, CardRepository インターフェース
│   ├── usecase/
│   │   ├── board.go          # BoardUseCase
│   │   └── card.go           # CardUseCase（move, archive含む）
│   ├── handler/
│   │   ├── board.go          # ボードCRUDハンドラ
│   │   ├── card.go           # カードCRUD + move + archive
│   │   └── ws.go             # WebSocketハンドラ
│   └── infra/
│       ├── yaml/
│       │   └── store.go      # BoardRepository, CardRepository の YAML実装
│       └── watcher/
│           └── watcher.go    # fsnotify監視 → WebSocket通知
├── web/                      # Reactフロントエンド
│   ├── src/
│   │   ├── App.tsx
│   │   ├── components/
│   │   │   ├── Board.tsx     # カンバンボード全体
│   │   │   ├── List.tsx      # リスト列
│   │   │   ├── Card.tsx      # カードコンポーネント
│   │   │   ├── CardModal.tsx # カード詳細・編集モーダル
│   │   │   └── ArchiveView.tsx # アーカイブ一覧
│   │   ├── hooks/
│   │   │   ├── useApi.ts     # REST API呼び出し
│   │   │   └── useWebSocket.ts # WebSocket接続
│   │   └── types.ts          # 型定義
│   ├── package.json
│   └── vite.config.ts
├── go.mod
├── go.sum
└── mise.toml                 # タスクランナー + ツールバージョン管理
```
