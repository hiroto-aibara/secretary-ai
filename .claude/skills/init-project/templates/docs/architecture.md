# アーキテクチャ

## 概要

<!-- プロジェクトの目的と概要を1-3文で記述 -->

## システム構成

<!-- システム全体の構成図をASCIIアートまたはMermaidで記述 -->

```
┌─────────────────┐            ┌──────────────────┐
│  Client         │── REST ──▶│  Server          │
└─────────────────┘            └──────────────────┘
```

## データ構造

<!-- データの保存形式・ディレクトリレイアウト・スキーマを記述 -->

### ディレクトリレイアウト

```
data/
└── ...
```

### スキーマ

<!-- 主要なデータ構造の定義を記述 -->

## バックエンドレイヤー設計

### レイヤー構成と依存方向

```
handler → usecase → domain ← infra
```

- `domain`: エンティティ + リポジトリインターフェース。他パッケージに依存しない
- `usecase`: ビジネスロジック。`domain` のインターフェースにのみ依存
- `handler`: HTTPリクエスト/レスポンス処理。`usecase` に依存
- `infra`: `domain` インターフェースの具体実装

### ルール

| # | ルール | 詳細 |
|---|--------|------|
| 1 | 依存方向は内側のみ | handler→usecase→domain の方向のみ許可 |
| 2 | インターフェース定義は domain | リポジトリ等のインターフェースは `domain` パッケージに定義 |
| 3 | DI はコンストラクタ注入 | フレームワーク不使用。`New*` 関数で依存を受け取る |
| 4 | DI 配線は main.go に集約 | `cmd/<app>/main.go` が唯一の配線ポイント |
| 5 | handler はロジックを持たない | リクエスト解析 + usecase 呼び出し + レスポンス構築のみ |
| 6 | usecase は infra を知らない | インターフェース経由でのみデータアクセス |

### インターフェース定義例

```go
// domain/repository.go
package domain

import "context"

type ExampleRepository interface {
    List(ctx context.Context) ([]Entity, error)
    Get(ctx context.Context, id string) (*Entity, error)
    Save(ctx context.Context, entity *Entity) error
    Delete(ctx context.Context, id string) error
}
```

### DI配線例（main.go）

```go
func main() {
    // infra
    store := infra.NewStore("./data")

    // usecase
    uc := usecase.NewExampleUseCase(store)

    // handler
    h := handler.NewExampleHandler(uc)

    // router
    r := chi.NewRouter()
    h.Register(r)

    http.ListenAndServe(":8080", r)
}
```

## 技術スタック

### バックエンド

| ライブラリ | 用途 |
|-----------|------|
| Go 1.22+  | 言語 |
<!-- 依存パッケージを列挙 -->

### フロントエンド

| ライブラリ | 用途 |
|-----------|------|
| React 18  | UI |
| TypeScript | 型安全性 |
| Vite      | ビルドツール |
<!-- 追加パッケージを列挙 -->

## ソースコード構造

```
project/
├── cmd/
│   └── <app>/
│       └── main.go
├── internal/
│   ├── domain/
│   ├── usecase/
│   ├── handler/
│   └── infra/
├── web/
│   └── src/
├── go.mod
└── mise.toml
```
