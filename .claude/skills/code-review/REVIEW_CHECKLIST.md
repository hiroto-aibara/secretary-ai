# OnboardAI コードレビューチェックリスト

このドキュメントは、OnboardAI プロジェクト固有のレビュー観点をまとめたものです。

## 1. アーキテクチャ / レイヤ構造

### Must（必須）

- [ ] **Domain層が Framework を import していない**
  - FastAPI, boto3, 外部SDK の import 禁止
  - 違反例: `from fastapi import HTTPException` in Domain
  - 正解: Domain は独自例外を定義（`app.domain.exceptions`）

- [ ] **レイヤ依存の方向が正しい**
  - API → Application → Domain
  - Infrastructure → Domain（DomainのIFに適合）
  - 違反: Domain → Infrastructure, Application → API

- [ ] **DI配線が Composition Root にある**
  - `app.bootstrap` または `main.py`（レイヤと同階層）で配線
  - 違反: `api` が直接 `infrastructure` を import

### Should（推奨）

- [ ] **1ファイル1責務**
  - 巨大なファイル（500行以上）は分割検討
  - 責務が混在していないか（API + ビジネスロジック等）

## 2. エラーハンドリング

### Must（必須）

- [ ] **各層で適切な例外を使用**
  - Domain: `InvalidCandidateDataError`（ビジネスルール違反）
  - Domain: `DriveServiceError`, `DocuSignServiceError`（外部サービスエラー）
  - Application: `OnboardingNotFoundError`（ユースケース失敗）
  - Infrastructure: Domain層の `ExternalServiceError` を使用（独自定義しない）
  - API: `@app.exception_handler` で統一変換

- [ ] **エラーコードが定義されている**
  - `DOMAIN.REASON` 形式（例: `ONBOARDING.CANDIDATE_INVALID`）
  - エラーコードは安定、メッセージは変更可

- [ ] **HTTPException が Domain/Application 層にない**
  - 違反例: Domain で `raise HTTPException(status_code=400)`
  - 正解: Domain で `raise InvalidCandidateDataError()` → API層で変換

### Should（推奨）

- [ ] **エラーメッセージが明確**
  - ユーザーが対応できる情報を含む
  - 内部情報（スタックトレース等）を露出しない

## 3. セキュリティ

### Must（必須）

- [ ] **認証情報がハードコードされていない**
  - `.env` は gitignore 済み
  - 本番は Secrets Manager 使用
  - 違反例: `API_KEY = "sk-abc123..."`

- [ ] **Webhook 署名検証がある**
  - DocuSign: `X-DocuSign-Signature-1` ヘッダー検証
  - 冪等性処理（重複イベント検出）
  - 違反例: 署名検証なしで payload 処理

- [ ] **入力検証がある**
  - API エンドポイントで型・値チェック
  - SQL injection 対策（パラメータ化クエリ）
  - XSS対策（エスケープ処理）

### Should（推奨）

- [ ] **エラーメッセージが情報を露出していない**
  - 違反例: `Database connection failed at host 192.168.1.1`
  - 正解: `Internal server error`（詳細はログのみ）

- [ ] **Rate limiting がある**
  - 認証エンドポイントで必須
  - 公開 API で推奨

## 4. テスト

### Must（必須）

- [ ] **Domain層のビジネスルールにテストがある**
  - 命名: `test_r<N>_<constraint>_<scenario>`
  - 例: `test_r1_1_cannot_send_before_generation()`

- [ ] **Infrastructure層の CRUD にテストがある**
  - create/read/update/delete の正常系・異常系
  - 楽観的ロックエラーのテスト
  - 命名: `test_<operation>_<expected_result>`

- [ ] **async テストに `@pytest.mark.asyncio` がある**
  - 違反例: `async def test_xxx()` のみ
  - 正解: `@pytest.mark.asyncio` デコレータ付与

### Should（推奨）

- [ ] **Fixture命名規則に従っている**
  - 共有fixture: `sample_*`（`tests/fixtures/domain.py`）
  - ビルダー: `*_builder`（`tests/<layer>/conftest.py`）

- [ ] **AAA パターンに従っている**
  - Arrange（準備）
  - Act（実行）
  - Assert（検証）

- [ ] **例外検証で型とメッセージを検証している**
  - `pytest.raises()` でメッセージも確認
  - 違反例: 型のみチェック

### Nice（改善）

- [ ] **テストカバレッジが十分**
  - Domain層: 90%以上
  - Infrastructure層: 70%以上

## 5. パフォーマンス

### Must（必須）

- [ ] **N+1 問題がない**
  - ループ内で DB/API 呼び出しがない
  - 違反例: `for item in items: db.get(item.id)`
  - 正解: バッチ取得または JOIN

### Should（推奨）

- [ ] **不要なループがない**
  - リスト内包表記や組み込み関数で代替可能か

- [ ] **キャッシュ可能なデータをキャッシュしている**
  - 頻繁にアクセスする静的データ

## 6. 可読性・保守性

### Should（推奨）

- [ ] **変数名が明確**
  - 違反例: `d`, `tmp`, `x`
  - 正解: `onboarding_data`, `candidate_email`

- [ ] **関数が単一責務**
  - 1関数50行以内を目安
  - 複雑な処理は分割

- [ ] **マジックナンバーがない**
  - 違反例: `if status == 3:`
  - 正解: `if status == DocumentStatus.COMPLETED:`

### Nice（改善）

- [ ] **Docstring がある**
  - 公開関数・メソッドに説明
  - 複雑なロジックに補足

- [ ] **型ヒントがある**
  - 公開関数・メソッドに型注釈
  - Python 3.11+ の新しい型構文を活用

## 7. Lambda / AWS 制約

### Must（必須）

- [ ] **タイムアウトを考慮している**
  - Lambda 最大15分
  - 長時間処理は分割またはStep Functions

- [ ] **冪等性がある**
  - Webhook/ジョブは重複実行を考慮
  - 条件付き更新でロック

### Should（推奨）

- [ ] **リトライ可能な設計**
  - 外部API呼び出しは Exponential Backoff
  - 部分失敗からの復旧考慮

## 8. ドキュメント

### Should（推奨）

- [ ] **破壊的変更時にドキュメント更新**
  - API変更 → `docs/mvp-design.md` 更新
  - デプロイ手順変更 → `CLAUDE.md` 更新

- [ ] **新機能追加時に実装計画更新**
  - `docs/implementation-plan.md` の該当タスクを完了済みに

### Nice（改善）

- [ ] **複雑なロジックにコメント**
  - WHYを説明（WHATはコードで自明）

## 9. Git / コミット

### Should（推奨）

- [ ] **コミットメッセージが明確**
  - 形式: `feat:`, `fix:`, `docs:`, `refactor:` 等
  - 1コミット1意図

- [ ] **無関係な変更が含まれていない**
  - フォーマット修正のみのコミットは別PR推奨

## 使い方

このチェックリストを元に、コードレビュー時に各項目を確認してください。特に **Must（必須）** 項目は、マージ前に必ず対応が必要です。
