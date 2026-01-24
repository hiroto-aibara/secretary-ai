# API仕様

## REST API

### エンドポイント一覧

| Method | Path | 説明 |
|--------|------|------|
<!-- エンドポイントを列挙 -->

### リクエスト/レスポンス例

<!-- 主要なエンドポイントのリクエスト/レスポンス例を記述 -->

## エラーレスポンス

全APIエンドポイントで統一されたエラー形式を使用する。

### 形式

```json
{
  "error": {
    "code": "not_found",
    "message": "Resource not found"
  }
}
```

### エラーコード一覧

| HTTPステータス | code | 説明 |
|---------------|------|------|
| 400 | `bad_request` | リクエストボディが不正 |
| 400 | `validation_error` | バリデーションエラー |
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
      {"field": "title", "message": "title is required"}
    ]
  }
}
```

## WebSocket

<!-- WebSocket を使用する場合に記述 -->

### 接続

```
ws://localhost:8080/ws
```

### メッセージ形式（サーバー → クライアント）

```json
{
  "type": "resource_updated",
  "resource_id": "example-id",
  "timestamp": "2026-01-01T00:00:00+09:00"
}
```

### イベントタイプ

| type | トリガー |
|------|---------|
<!-- イベントタイプを列挙 -->
