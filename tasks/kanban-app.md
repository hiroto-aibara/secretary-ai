# Task: kanban-app

## 概要

YAML駆動のカンバン式タスク管理アプリケーションのフルスタック実装。
GoバックエンドでYAMLファイルのCRUD + ファイル監視 + WebSocket通知を行い、
ReactフロントエンドでカンバンUI（ドラッグ&ドロップ、アーカイブ）を提供する。

## 要件

### バックエンド（Go）

#### domain レイヤー
- [ ] `internal/domain/board.go` — Board, List エンティティ定義
- [ ] `internal/domain/card.go` — Card エンティティ定義（ID, title, list, order, description, labels, archived, timestamps）
- [ ] `internal/domain/repository.go` — BoardRepository, CardRepository インターフェース定義
- [ ] `internal/domain/errors.go` — ErrNotFound, ErrValidation, ErrConflict 定義

#### usecase レイヤー
- [ ] `internal/usecase/board.go` — BoardUseCase（List, Get, Create, Update, Delete）
- [ ] `internal/usecase/card.go` — CardUseCase（List, Get, Create, Update, Delete, Move, Archive）

#### infra レイヤー
- [ ] `internal/infra/yaml/store.go` — BoardRepository, CardRepository の YAML実装
  - `.tasks/boards/<id>/board.yaml` の読み書き
  - `.tasks/boards/<id>/cards/<cardId>.yaml` の読み書き
  - カードID生成（日付+連番: `YYYYMMDD-NNN`）
- [ ] `internal/infra/watcher/watcher.go` — fsnotify で `.tasks/` を監視
  - ファイル変更検知 → WebSocket Hub に通知
  - board.yaml 変更 → `board_updated` イベント
  - card YAML 変更 → `card_updated` イベント

#### handler レイヤー
- [ ] `internal/handler/board.go` — ボード CRUD ハンドラ
  - GET /api/boards, POST /api/boards
  - GET /api/boards/:id, PUT /api/boards/:id, DELETE /api/boards/:id
- [ ] `internal/handler/card.go` — カード CRUD + move + archive ハンドラ
  - GET /api/boards/:id/cards(?archived=true)
  - POST /api/boards/:id/cards
  - GET/PUT/DELETE /api/boards/:id/cards/:cardId
  - PATCH /api/boards/:id/cards/:cardId/move
  - PATCH /api/boards/:id/cards/:cardId/archive
- [ ] `internal/handler/ws.go` — WebSocket ハンドラ（Hub管理、接続/切断/ブロードキャスト）
- [ ] `internal/handler/response.go` — 共通レスポンスヘルパー（respondJSON, writeError）

#### エントリポイント
- [ ] `cmd/taskmgr/main.go` — DI配線 + chi ルーター設定 + サーバー起動
  - フロントエンド静的ファイル配信（embed.FS）
  - Graceful shutdown

### フロントエンド（React + TypeScript）

#### 型定義・API
- [ ] `web/src/types.ts` — Board, List, Card, APIレスポンス型
- [ ] `web/src/hooks/useApi.ts` — REST API 呼び出し（fetch ベース）
- [ ] `web/src/hooks/useWebSocket.ts` — WebSocket 接続管理 + 自動再接続

#### コンポーネント
- [ ] `web/src/components/Board.tsx` — カンバンボード全体（リスト列の横並び）
- [ ] `web/src/components/List.tsx` — リスト列（カード一覧 + 新規カード追加）
- [ ] `web/src/components/Card.tsx` — カードコンポーネント（タイトル、ラベル表示）
- [ ] `web/src/components/CardModal.tsx` — カード詳細・編集モーダル
- [ ] `web/src/components/ArchiveView.tsx` — アーカイブ一覧 + 復元操作
- [ ] `web/src/App.tsx` — ボード選択 + メインレイアウト

#### ドラッグ&ドロップ
- [ ] @dnd-kit/core + @dnd-kit/sortable でカード移動
- [ ] リスト間移動（list変更）+ リスト内並び替え（order変更）
- [ ] ドロップ時に PATCH /move API を呼び出し

#### リアルタイム更新
- [ ] WebSocket でイベント受信 → 該当ボード/カードのデータ再取得
- [ ] 自動再接続（切断時）

### ビルド統合
- [ ] `web/dist/` を Go バイナリに `embed.FS` で埋め込み
- [ ] `mise run build` でフロント→バックエンド一括ビルド
- [ ] 単一バイナリ `bin/taskmgr` でフロント配信 + API提供

## 技術仕様

- 使用技術:
  - Backend: Go 1.22+, chi/v5, yaml.v3, fsnotify, gorilla/websocket
  - Frontend: React 18, TypeScript, Vite, @dnd-kit, CSS Modules
- レイヤー設計: `handler → usecase → domain ← infra`
- データ保存: `.tasks/boards/<id>/` 配下のYAMLファイル
- API: REST（JSON）+ WebSocket（通知のみ）

## 受け入れ基準

- [ ] `mise run build && ./bin/taskmgr` でサーバー起動し、ブラウザでカンバンUIが表示される
- [ ] REST API でボード・カードのCRUD操作ができ、YAMLファイルが正しく更新される
- [ ] Web UIでカードのドラッグ&ドロップ移動が動作する
- [ ] カードのアーカイブ・復元が動作する
- [ ] YAMLファイルを直接編集した場合、WebSocket経由でUIがリアルタイム更新される
- [ ] `mise run lint` がエラーなしで通る
- [ ] `mise run test` が全テストパスする
- [ ] golangci-lint の depguard ルール（レイヤー依存方向）に違反がない

## 品質基準

- テストカバレッジ: 全体80%以上（domain/usecase: 90%, handler: 80%, infra: 70%）
- エラーハンドリング: 統一エラーレスポンス形式（`{error: {code, message}}`）
- ログ: log/slog 使用、構造化ログ
- セキュリティ: 入力バリデーション（title必須、list存在チェック等）

## 参考リソース

- docs/architecture.md — レイヤー設計・依存方向ルール・DI配線例
- docs/api.md — REST API 仕様・エラーレスポンス・WebSocket仕様
- docs/conventions.md — エラーハンドリング・テスト・ログ規約
