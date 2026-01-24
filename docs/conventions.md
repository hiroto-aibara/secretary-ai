# 開発規約

## エラーハンドリング

### 方針

- 独自エラー型を `domain` パッケージに定義
- `errors.Is` / `errors.As` で判別可能にする
- 各レイヤーの責務を明確に分離

### レイヤーごとの責務

| レイヤー | 責務 |
|----------|------|
| domain | エラー型の定義 |
| infra | OS/IOエラーをdomainエラーにラップして返却 |
| usecase | domainエラーをそのまま返却（追加ラップ不要） |
| handler | domainエラーをHTTPステータス + エラーレスポンスに変換 |

### エラー型定義

```go
// domain/errors.go
package domain

import "fmt"

type ErrNotFound struct {
    Resource string
    ID       string
}

func (e *ErrNotFound) Error() string {
    return fmt.Sprintf("%s %s not found", e.Resource, e.ID)
}

type ErrValidation struct {
    Field   string
    Message string
}

func (e *ErrValidation) Error() string {
    return fmt.Sprintf("validation error: %s %s", e.Field, e.Message)
}

type ErrConflict struct {
    Resource string
    ID       string
}

func (e *ErrConflict) Error() string {
    return fmt.Sprintf("%s %s already exists", e.Resource, e.ID)
}
```

### handler でのマッピング

```go
func writeError(w http.ResponseWriter, err error) {
    var notFound *domain.ErrNotFound
    var validation *domain.ErrValidation
    var conflict *domain.ErrConflict

    switch {
    case errors.As(err, &notFound):
        respondJSON(w, http.StatusNotFound, errorResponse("not_found", err.Error()))
    case errors.As(err, &validation):
        respondJSON(w, http.StatusBadRequest, errorResponse("validation_error", err.Error()))
    case errors.As(err, &conflict):
        respondJSON(w, http.StatusConflict, errorResponse("conflict", err.Error()))
    default:
        slog.Error("unexpected error", "error", err)
        respondJSON(w, http.StatusInternalServerError, errorResponse("internal_error", "internal server error"))
    }
}
```

### 禁止事項

- `panic` をエラーハンドリングに使わない（プログラムの不整合のみ）
- エラーメッセージにスタックトレースを含めない（ログに出力）
- ユーザー向けレスポンスに内部実装の詳細を含めない

## テスト方針

### テスト対象と方法

| レイヤー | テスト対象 | 方法 |
|----------|-----------|------|
| domain | エンティティのバリデーション・振る舞い | ユニットテスト |
| usecase | ビジネスロジック | リポジトリインターフェースをモック |
| infra/yaml | YAML読み書きの正確性 | 一時ディレクトリ使用の統合テスト |
| handler | HTTPリクエスト/レスポンス | httptest + usecaseモック |

### テストスタイル

- **テーブル駆動テスト**を標準とする

```go
func TestCardUseCase_Archive(t *testing.T) {
    tests := []struct {
        name    string
        cardID  string
        setup   func(*mockCardRepo)
        wantErr error
    }{
        {
            name:   "success",
            cardID: "20260124-001",
            setup: func(m *mockCardRepo) {
                m.card = &domain.Card{ID: "20260124-001", Archived: false}
            },
            wantErr: nil,
        },
        {
            name:   "not found",
            cardID: "nonexistent",
            setup: func(m *mockCardRepo) {
                m.err = &domain.ErrNotFound{Resource: "card", ID: "nonexistent"}
            },
            wantErr: &domain.ErrNotFound{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mock := &mockCardRepo{}
            tt.setup(mock)
            uc := usecase.NewCardUseCase(mock, nil)

            _, err := uc.Archive(context.Background(), "board-1", tt.cardID)

            if tt.wantErr != nil {
                if !errors.As(err, &tt.wantErr) {
                    t.Errorf("got %v, want %v", err, tt.wantErr)
                }
            } else if err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

### モック

- 手書きモック（インターフェースが小さいため外部ライブラリ不要）
- テストファイル内に定義（`*_test.go`）

```go
type mockCardRepo struct {
    card *domain.Card
    err  error
}

func (m *mockCardRepo) Get(ctx context.Context, boardID, cardID string) (*domain.Card, error) {
    return m.card, m.err
}
// ... 他メソッドも同様
```

### テスト命名規則

```
TestXxx_Method            # 正常系
TestXxx_Method_異常ケース  # 異常系（テーブル内で分岐）
```

### 実行

```bash
mise exec -- go test ./internal/...
```

## ログ

### ライブラリ

`log/slog`（Go標準ライブラリ）を使用。外部依存なし。

### ログレベル

| レベル | 用途 | 例 |
|--------|------|-----|
| Error | 処理継続不能な障害 | YAML書き込み失敗、ファイル監視エラー |
| Warn | 回復可能だが注意が必要 | 不正なYAMLファイルのスキップ |
| Info | 操作の記録 | サーバー起動、カード作成/移動 |
| Debug | 開発時の詳細情報 | リクエストボディ、ファイル変更検知 |

### 書式

```go
// 構造化ログを使用（キーバリューペア）
slog.Info("card created", "board_id", boardID, "card_id", card.ID)
slog.Error("failed to save card", "board_id", boardID, "card_id", cardID, "error", err)

// メッセージは小文字開始、簡潔に
// Good
slog.Info("card moved", "card_id", id, "from", oldList, "to", newList)
// Bad
slog.Info("The card has been successfully moved to a new list", ...)
```

### レイヤーごとのルール

| レイヤー | ログ |
|----------|------|
| handler | リクエスト受信（Info）、エラーレスポンス（Warn/Error） |
| usecase | 書かない（handler に委譲） |
| infra | IO障害（Error）、リトライ（Warn） |
| watcher | ファイル変更検知（Debug）、監視エラー（Error） |

### 禁止事項

- `fmt.Println` をログに使わない
- エラーを握りつぶさない（ログ出力するか、上位に返す）
- パスワード・トークン等の機密情報をログに含めない
