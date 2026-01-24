---
name: create-task
description: Creates a TASK.md file based on user's description. Use this before creating worktrees to define task requirements and quality standards.
allowed-tools: Write, Read, Bash(mkdir:*)
---

# Task Creator

ユーザーの説明を元にTASK.mdファイルを生成します。

## 概要

このSkillは以下を実行します：

1. ユーザーからタスクの説明を受け取る
2. Claude Codeが要件・受け入れ基準・品質基準を含むTASK.mdを生成
3. `tasks/<feature-name>.md` として保存

## 使用方法

```bash
/create-task <feature-name>
```

### 例

```bash
/create-task user-auth
# ユーザー: "JWTを使った認証機能。ログイン、ログアウト、トークンリフレッシュが必要"
# → tasks/user-auth.md が生成される
```

## 生成されるファイル

```
tasks/
├── user-auth.md
├── dashboard.md
└── api-v2.md
```

## TASK.mdの構成

生成されるTASK.mdには以下のセクションが含まれます：

```markdown
# Task: <feature-name>

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
```

## 詳細

詳細については [REFERENCE.md](REFERENCE.md) を参照してください。
