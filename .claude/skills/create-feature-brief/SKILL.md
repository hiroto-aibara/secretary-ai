---
name: create-feature-brief
description: Creates a Feature Brief document defining the purpose, scope, and use cases of a new feature. Use this as the first step before design docs and task breakdown.
allowed-tools: Write, Read, Bash(mkdir:*), AskUserQuestion
---

# Feature Brief Creator

ユーザーの説明を元に Feature Brief ドキュメントを生成します。

## 概要

このSkillは以下を実行します：

1. ユーザーから機能の目的・スコープを聞き取る
2. 目的・ゴール・スコープ・ユースケースを含む Feature Brief を生成
3. `docs/<feature-name>-brief.md` として保存

## 使用方法

```bash
/create-feature-brief <feature-name>
```

### 例

```bash
/create-feature-brief user-auth
# ユーザー: "JWTを使った認証機能。ログイン、ログアウト、トークンリフレッシュが必要"
# → docs/user-auth-brief.md が生成される
```

## Feature Brief の位置づけ

```
Feature Brief（なぜ・何を）  ← このスキルで作成
    ↓
Design Doc（どうやって）     ← /create-design-doc で作成
    ↓
Task File（実装指示）        ← /create-task で作成
```

## 実行手順

1. **要件ヒアリング**
   - 機能の目的・背景を確認
   - 対象ユーザー/アクターを確認
   - スコープ（やること/やらないこと）を明確化
   - 不明点があれば質問

2. **Feature Brief 生成**
   - [TEMPLATE.md](TEMPLATE.md) のフォーマットに従って生成
   - ユースケースは具体的なフローで記述
   - 非機能要件（セキュリティ、パフォーマンス等）も含める
   - オープンクエスチョンがあれば明記

3. **ユーザーレビュー**
   - 生成したドキュメントを表示
   - 修正リクエストがあれば反映

## 生成ファイル

```
docs/
├── user-auth-brief.md
├── notification-brief.md
└── ...
```

## テンプレート

[TEMPLATE.md](TEMPLATE.md) にフォーマットがあります。

## 関連スキル

- `/create-design-doc` - Design Doc 作成（Feature Brief 承認後）
- `/create-task` - タスクファイル作成（Design Doc 作成後）
