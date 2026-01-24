---
name: create-vk-task
description: Creates a task in vibe-kanban via MCP. Use this to register tasks for tracking and workspace management.
allowed-tools: mcp__vibe_kanban__list_projects, mcp__vibe_kanban__create_task, mcp__vibe_kanban__list_tasks
---

# Vibe-Kanban Task Creator

ユーザーの説明を元にvibe-kanbanにタスクを登録します。

## 概要

このSkillは以下を実行します：

1. ユーザーからタスクの説明を受け取る
2. Claude Codeが要件・受け入れ基準・品質基準を含むタスク説明を生成
3. vibe-kanban MCP経由でタスクを登録

## 使用方法

```bash
/create-vk-task <タスクタイトル>
```

### 例

```bash
/create-vk-task JWTを使った認証機能
# ユーザー: "ログイン、ログアウト、トークンリフレッシュが必要"
# → vibe-kanbanにタスクが登録される
```

## プロジェクト設定

| プロジェクト | ID |
|-------------|-----|
| OnboardAI | `dda09758-b4ed-4960-be35-a00157366e50` |

## タスク説明のフォーマット

タスクの説明（description）には以下のセクションを含めます：

```markdown
## 概要
タスクの概要説明

## 要件
- [ ] 要件1
- [ ] 要件2
- [ ] 要件3

## 技術仕様
- 使用技術:
- 関連ファイル:
- API エンドポイント:

## 受け入れ基準
- [ ] 全テストがパス
- [ ] ESLint/Prettier エラーなし
- [ ] コードレビュー完了

## 品質基準
- テストカバレッジ: 80%以上
- パフォーマンス: レスポンス200ms以内
- セキュリティ: 入力バリデーション必須

## 参考リソース
- 関連Issue:
- デザイン:
- API仕様:

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

## 実行手順

1. **プロジェクト確認**
   - `mcp__vibe_kanban__list_projects` でプロジェクト一覧取得
   - OnboardAIプロジェクト（`dda09758-b4ed-4960-be35-a00157366e50`）を使用

2. **タスク情報の収集**
   - ユーザーからタスクの詳細を聞き取る
   - 不明点があれば質問して明確化

3. **タスク説明の生成**
   - 上記フォーマットに従って説明文を生成
   - 要件は具体的かつ検証可能な形で記述
   - 「完了時の手順」セクションは必ず含める

4. **タスク登録**
   - `mcp__vibe_kanban__create_task` でタスク作成
   - 必須パラメータ:
     - `project_id`: プロジェクトID
     - `title`: タスクタイトル
     - `description`: 生成した説明文

5. **結果の表示**
   - 作成されたタスクIDを表示
   - 次のステップ（ワークスペース作成）を案内

## 命名規則

| タスク種別 | タイトル例 |
|-----------|-----------|
| 新機能 | `feature: ユーザー認証機能` |
| バグ修正 | `fix: ログインエラーの修正` |
| リファクタリング | `refactor: APIクライアントの整理` |
| パフォーマンス改善 | `perf: クエリ最適化` |

## 関連スキル

- `/code-review` - コードレビュー（セルフレビュー用）
- `/create-pr` - PR作成とワークツリー管理

## 詳細

詳細については [REFERENCE.md](REFERENCE.md) を参照してください。
