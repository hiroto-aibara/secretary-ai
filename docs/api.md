# API仕様

## REST API

### エンドポイント一覧

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

#### POST /api/boards

```json
// Request
{
  "id": "my-project",
  "name": "My Project",
  "lists": [
    {"id": "todo", "name": "Todo"},
    {"id": "in-progress", "name": "In Progress"},
    {"id": "done", "name": "Done"}
  ]
}

// Response 201
{
  "id": "my-project",
  "name": "My Project",
  "lists": [...]
}
```

#### GET /api/boards/:id

```json
// Response 200
{
  "id": "my-project",
  "name": "My Project",
  "lists": [
    {"id": "todo", "name": "Todo"},
    {"id": "in-progress", "name": "In Progress"},
    {"id": "done", "name": "Done"}
  ]
}
```

#### POST /api/boards/:id/cards

```json
// Request
{
  "title": "新機能の設計",
  "list": "todo",
  "description": "詳細な仕様を決める",
  "labels": ["feature"]
}

// Response 201
{
  "id": "20260124-001",
  "title": "新機能の設計",
  "list": "todo",
  "order": 0,
  "description": "詳細な仕様を決める",
  "labels": ["feature"],
  "archived": false,
  "created_at": "2026-01-24T10:00:00+09:00",
  "updated_at": "2026-01-24T10:00:00+09:00"
}
```

#### PUT /api/boards/:id/cards/:cardId

```json
// Request（部分更新可）
{
  "title": "更新されたタイトル",
  "description": "更新された説明"
}

// Response 200
{
  "id": "20260124-001",
  "title": "更新されたタイトル",
  ...
}
```

#### PATCH /api/boards/:id/cards/:cardId/move

```json
// Request
{
  "list": "in-progress",
  "order": 2
}

// Response 200
{
  "id": "20260124-001",
  "list": "in-progress",
  "order": 2,
  ...
}
```

#### PATCH /api/boards/:id/cards/:cardId/archive

```json
// Request
{
  "archived": true
}

// Response 200
{
  "id": "20260124-001",
  "archived": true,
  ...
}
```

#### GET /api/boards/:id/cards?archived=true

アーカイブ済みカードを含む全カードを返却。
`archived` パラメータ省略時はアクティブカードのみ。

## エラーレスポンス

全APIエンドポイントで統一されたエラー形式を使用する。

### 形式

```json
{
  "error": {
    "code": "not_found",
    "message": "Card 20260124-001 not found in board my-project"
  }
}
```

### エラーコード一覧

| HTTPステータス | code | 説明 |
|---------------|------|------|
| 400 | `bad_request` | リクエストボディが不正 |
| 400 | `validation_error` | バリデーションエラー（必須フィールド不足等） |
| 404 | `not_found` | リソースが存在しない |
| 409 | `conflict` | IDの重複等 |
| 500 | `internal_error` | サーバー内部エラー |

### バリデーションエラーの詳細

```json
{
  "error": {
    "code": "validation_error",
    "message": "Validation failed",
    "details": [
      {"field": "title", "message": "title is required"},
      {"field": "list", "message": "list 'invalid-list' does not exist in board"}
    ]
  }
}
```

## WebSocket

### 接続

```
ws://localhost:8080/ws
```

### メッセージ形式（サーバー → クライアント）

ファイル変更検知時に以下のイベントをブロードキャスト:

```json
{
  "type": "board_updated",
  "board_id": "my-project",
  "timestamp": "2026-01-24T15:00:00+09:00"
}
```

### イベントタイプ

| type | トリガー |
|------|---------|
| `board_updated` | board.yaml の変更 |
| `card_updated` | カードYAMLの作成・変更・削除 |

クライアントはイベント受信後、必要なAPIを再呼び出ししてデータを最新化する。
（差分配信ではなく、通知のみを行うシンプルな設計）
