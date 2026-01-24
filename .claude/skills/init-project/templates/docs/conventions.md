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

type ErrConflict struct {
    Resource string
    ID       string
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
| infra | データ読み書きの正確性 | 一時ディレクトリ使用の統合テスト |
| handler | HTTPリクエスト/レスポンス | httptest + usecaseモック |

### テストスタイル

- **テーブル駆動テスト**を標準とする

```go
func TestExample_Method(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr error
    }{
        {
            name:  "success",
            input: "valid",
            want:  "expected",
        },
        {
            name:    "not found",
            input:   "invalid",
            wantErr: &domain.ErrNotFound{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### モック

- 手書きモック（インターフェースが小さいため外部ライブラリ不要）
- テストファイル内に定義（`*_test.go`）

### テスト命名規則

```
TestXxx_Method            # 正常系
TestXxx_Method_異常ケース  # 異常系（テーブル内で分岐）
```

### カバレッジ目標

| レイヤー | 目標 | 備考 |
|----------|------|------|
| domain | 90% | 純粋ロジック、テスト容易 |
| usecase | 90% | ビジネスロジック中核 |
| handler | 80% | httptest で網羅 |
| infra | 70% | I/O依存、統合テスト中心 |
| **全体** | **80%** | CI で閾値チェック |

### 実行

```bash
mise run test             # テスト実行（カバレッジ表示付き）
mise run test:coverage    # カバレッジレポート生成 + 閾値チェック
```

## ログ

### ライブラリ

`log/slog`（Go標準ライブラリ）を使用。

### ログレベル

| レベル | 用途 | 例 |
|--------|------|-----|
| Error | 処理継続不能な障害 | 書き込み失敗、監視エラー |
| Warn | 回復可能だが注意が必要 | 不正データのスキップ |
| Info | 操作の記録 | サーバー起動、リソース作成/更新 |
| Debug | 開発時の詳細情報 | リクエストボディ、変更検知 |

### 書式

```go
// 構造化ログを使用（キーバリューペア）
slog.Info("resource created", "id", id)
slog.Error("failed to save", "id", id, "error", err)

// メッセージは小文字開始、簡潔に
```

### レイヤーごとのルール

| レイヤー | ログ |
|----------|------|
| handler | リクエスト受信（Info）、エラーレスポンス（Warn/Error） |
| usecase | 書かない（handler に委譲） |
| infra | IO障害（Error）、リトライ（Warn） |

### 禁止事項

- `fmt.Println` をログに使わない
- エラーを握りつぶさない（ログ出力するか、上位に返す）
- 機密情報をログに含めない
