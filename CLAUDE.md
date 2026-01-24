# プロジェクト規約

## 必読ドキュメント

- docs/architecture.md — レイヤー設計・依存方向ルール
- docs/conventions.md — エラーハンドリング・テスト・ログ規約
- docs/api.md — API仕様

## レイヤー依存ルール（最重要）

```
handler → usecase → domain ← infra
```

- domain は他の internal パッケージを import しない
- usecase は domain のみ import 可（handler, infra 禁止）
- handler は usecase, domain のみ import 可（infra 禁止）
- infra は domain のみ import 可

## コマンド

```bash
mise run fmt              # フォーマット（goimports + prettier）
mise run lint             # リントチェック（golangci-lint + tsc + eslint）
mise run test             # テスト実行（カバレッジ表示付き）
mise run test:coverage    # カバレッジレポート + 閾値チェック（80%）
mise run dev              # Go サーバー起動
mise run dev:front        # フロント dev サーバー起動
mise run build            # プロダクションビルド
```

## コード規約の要点

- エラー型は `internal/domain/errors.go` に定義
- エラー判別は `errors.As` / `errors.Is` を使用
- ログは `log/slog` を使用（構造化ログ、小文字開始）
- テストはテーブル駆動、モックは手書き、カバレッジ目標80%
- DI はコンストラクタ注入、配線は `cmd/taskmgr/main.go` に集約
