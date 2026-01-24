# タスク管理ツール設計書

## 概要

YAML駆動のカンバン式タスク管理ツール。
タスクの実態はYAMLファイルで、Claude Code（CLI）による直接編集と、Web UIからのREST API経由操作の両方に対応する。

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

## アーキテクチャ

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

## REST API

| Method | Path | 説明 |
|--------|------|------|
| GET    | `/api/boards` | ボード一覧 |
| POST   | `/api/boards` | ボード作成 |
| GET    | `/api/boards/:id` | ボード詳細（リスト情報含む） |
| PUT    | `/api/boards/:id` | ボード更新（リスト追加・名前変更等） |
| DELETE | `/api/boards/:id` | ボード削除 |
| GET    | `/api/boards/:id/cards` | カード一覧（`?archived=true`でアーカイブ含む） |
| POST   | `/api/boards/:id/cards` | カード作成 |
| GET    | `/api/boards/:id/cards/:cardId` | カード詳細 |
| PUT    | `/api/boards/:id/cards/:cardId` | カード更新 |
| DELETE | `/api/boards/:id/cards/:cardId` | カード削除 |
| PATCH  | `/api/boards/:id/cards/:cardId/move` | カード移動（list, order変更） |
| PATCH  | `/api/boards/:id/cards/:cardId/archive` | アーカイブ/復元トグル |

### リクエスト/レスポンス例

#### POST /api/boards/:id/cards

```json
{
  "title": "新機能の設計",
  "list": "todo",
  "description": "詳細な仕様を決める",
  "labels": ["feature"]
}
```

#### PATCH /api/boards/:id/cards/:cardId/move

```json
{
  "list": "in-progress",
  "order": 2
}
```

#### PATCH /api/boards/:id/cards/:cardId/archive

```json
{
  "archived": true
}
```

## アーカイブ仕様

- カードYAMLの `archived: true` フラグで管理
- アーカイブされたカードはカンバン上のリストに表示されない
- Web UIにアーカイブ一覧ビューを用意（フィルタ切り替え）
- 復元操作で `archived: false` に戻し、元のリストに復帰

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
│       └── main.go           # エントリポイント（サーバー起動）
├── internal/
│   ├── handler/
│   │   ├── board.go          # ボードCRUDハンドラ
│   │   ├── card.go           # カードCRUD + move + archive
│   │   └── ws.go             # WebSocketハンドラ
│   ├── model/
│   │   ├── board.go          # Board, List構造体
│   │   └── card.go           # Card構造体
│   ├── store/
│   │   └── yaml_store.go     # YAMLファイル読み書き操作
│   └── watcher/
│       └── watcher.go        # ファイル監視 → WebSocket通知
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
└── Makefile
```

## 実装ステップ

### Step 1: プロジェクト初期化
- Go module初期化、依存パッケージ追加
- Makefile作成（build, dev, cleanターゲット）

### Step 2: モデル定義
- Board, List, Card構造体
- YAMLタグ + JSONタグ付与

### Step 3: YAMLストア
- ボード一覧取得（ディレクトリスキャン）
- ボード読み込み（board.yaml解析）
- カードCRUD（YAML読み書き）
- カードID生成（日付+連番: `YYYYMMDD-NNN`）

### Step 4: HTTPハンドラ
- chiルーター設定
- ボードAPI実装
- カードAPI実装（move, archive含む）
- CORSミドルウェア（開発時用）

### Step 5: ファイルウォッチャー + WebSocket
- fsnotifyで`.tasks/`配下を再帰監視
- 変更検知時にWebSocketで接続クライアントに通知
- デバウンス処理（短時間の連続変更をまとめる）

### Step 6: サーバーエントリポイント
- ルーティング統合
- フロントエンドの静的ファイル配信（`embed.FS`）
- graceful shutdown

### Step 7: フロントエンド実装
- Vite + React + TypeScript初期化
- カンバンボードUI（リスト列 + カード）
- dnd-kitによるドラッグ&ドロップ
- カード作成・編集モーダル
- アーカイブ一覧ビュー
- WebSocketでリアルタイム更新

### Step 8: ビルド統合
- `make build`: フロントビルド → Goバイナリに組み込み
- `make dev`: フロント(Vite dev) + Go(air)を並行起動

## 検証方法

1. `make build && ./taskmgr` でサーバー起動
2. ブラウザで `http://localhost:8080` にアクセスしカンバン表示確認
3. REST APIをcurlで操作 → YAMLファイルが更新されることを確認
4. YAMLファイルを直接編集 → WebSocket経由でUIがリアルタイム更新されることを確認
5. Web UIでカードのD&D移動・アーカイブ・復元が動作することを確認
6. 複数ボードの切り替えが正しく動作することを確認
