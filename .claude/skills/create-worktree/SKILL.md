---
name: create-worktree
description: Creates a git worktree for parallel feature development. Use after planning to prepare an isolated development environment with all necessary environment files.
allowed-tools: Bash(git:*), Bash(mkdir:*), Bash(cp:*), Bash(chmod:*), Bash(bash:*), Bash(make:*)
---

# Git Worktree Creator

planモード終了後、feature開発用の独立したworktree環境を自動作成します。

## 概要

このSkillは以下を自動で実行します：

1. `.worktrees/<feature-name>/` ディレクトリにworktreeを作成
2. `feature/<feature-name>` ブランチを新規作成
3. 環境変数ファイル（`.env`, `.envrc` など）を自動コピー
4. `make setup` で開発環境をセットアップ

## 使用方法

### 基本的な使い方

```bash
# スクリプトを実行
bash .claude/skills/create-worktree/scripts/create_worktree.sh <feature-name>

# 例: user-auth 機能を開発する場合
bash .claude/skills/create-worktree/scripts/create_worktree.sh user-auth
```

### 実行結果

```
.worktrees/user-auth/     # worktreeディレクトリ
├── .env                  # ルートからコピー
├── .envrc                # ルートからコピー
├── modules/
│   ├── frontend/.env*    # frontendの環境変数
│   ├── backend/.env      # backendの環境変数
│   └── agent/.env        # agentの環境変数（あれば）
└── ...（その他のファイル）
```

## コピーされる環境変数ファイル

- ルート: `.env`, `.envrc`
- frontend: `.env`, `.env.local`, `.env.dev`, `.env.prd`, `.env.test`
- backend: `.env`
- agent: `.env`（存在する場合）

## 作業完了後

### 1. PR作成

```bash
/create-pr
```

詳細は [create-pr スキル](../create-pr/SKILL.md) を参照してください。

### 2. worktree削除（PRマージ後）

```bash
/cleanup-worktree
```

詳細は [cleanup-worktree スキル](../cleanup-worktree/SKILL.md) を参照してください。

### 手動でworktreeを削除する場合

```bash
git worktree remove .worktrees/<feature-name>
```

## 詳細

詳細については [REFERENCE.md](REFERENCE.md) を参照してください。
