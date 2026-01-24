# Vibe-Kanban Task Creator - Reference

## vibe-kanbanとの統合

このスキルはvibe-kanban MCPサーバーと連携してタスク管理を行います。

### ワークフロー全体像

```
1. /create-vk-task でタスク登録
   ↓
2. start_workspace_session でワークスペース作成
   - git worktree + ブランチが自動作成
   - mise run setup が自動実行
   ↓
3. タスク実装
   ↓
4. /code-review でセルフレビュー
   - 指摘事項を修正
   ↓
5. /create-pr でPR作成
   ↓
6. タスクステータスを inreview に更新
```

## タスク説明の書き方

### 良い例

```markdown
## 概要
DocuSign署名完了後のWebhookを受け取り、署名済みPDFをGoogle Driveに保存する機能を実装する。

## 要件
- [ ] POST /api/webhooks/docusign エンドポイント作成
- [ ] DocuSign署名検証ミドルウェア実装
- [ ] 署名済みPDF取得処理
- [ ] Google Drive保存処理
- [ ] 冪等性の担保（同一envelope_idの重複処理防止）
- [ ] 監査ログ記録

## 技術仕様
- 使用技術: FastAPI, boto3, google-api-python-client
- 関連ファイル:
  - app/api/routes/webhooks.py
  - app/application/usecases/docusign_webhook.py
  - app/infrastructure/docusign/
- API エンドポイント: POST /api/webhooks/docusign

## 受け入れ基準
- [ ] Webhook受信〜Drive保存のユニットテスト
- [ ] 冪等性テスト（同一リクエスト2回で1回のみ処理）
- [ ] 署名検証失敗時の401レスポンステスト
- [ ] ESLint/Prettier エラーなし
- [ ] mypy型チェックパス

## 品質基準
- テストカバレッジ: 90%以上（重要な外部連携のため）
- パフォーマンス: Webhook応答 3秒以内
- セキュリティ:
  - 署名検証必須
  - エラー時に機密情報を露出しない

## 参考リソース
- Issue #15: DocuSign Webhook統合
- docs/mvp-design.md: MVP設計仕様
- DocuSign Connect Guide

## 完了時の手順

タスク完了後、以下の手順を実行してください：

1. **セルフレビュー**
   `/code-review` スキルを使用してコードレビューを実施
   - 指摘事項があれば修正
   - 全ての指摘が解消されるまで繰り返す

2. **PR作成**
   `/create-pr` スキルを使用してPull Requestを作成
   - PRの説明にはタスクの概要と変更内容を記載
   - レビュワーを指定（必要に応じて）

3. **タスクステータス更新**
   PRが作成されたらタスクのステータスを `inreview` に更新
```

### 悪い例

```markdown
## 概要
Webhook作る

## 要件
- [ ] Webhookを実装
```

→ 具体性がなく、完了条件が不明確

## MCP ツール

### list_projects

プロジェクト一覧を取得します。

```
mcp__vibe_kanban__list_projects
```

### create_task

タスクを作成します。

```
mcp__vibe_kanban__create_task
  project_id: "dda09758-b4ed-4960-be35-a00157366e50"
  title: "feature: DocuSign Webhook統合"
  description: "## 概要\n..."
```

### list_tasks

タスク一覧を取得します。

```
mcp__vibe_kanban__list_tasks
  project_id: "dda09758-b4ed-4960-be35-a00157366e50"
```

## タスクのライフサイクル

| ステータス | 説明 |
|-----------|------|
| `todo` | 作成直後、未着手 |
| `inprogress` | ワークスペース作成後、作業中 |
| `inreview` | PR作成後、レビュー中 |
| `done` | PRマージ後、完了 |
| `cancelled` | キャンセル |

## 関連ドキュメント

- [CLAUDE.md](../../../CLAUDE.md) - プロジェクト全体のガイド
- [docs/mvp-design.md](../../../docs/mvp-design.md) - MVP設計仕様
- [docs/dev-rules.md](../../../docs/dev-rules.md) - 開発規約
