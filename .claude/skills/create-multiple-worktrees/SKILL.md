---
name: create-multiple-worktrees
description: Creates git worktrees from TASK.md files for parallel feature development. Each worktree gets its own TASK.md with requirements and quality standards.
allowed-tools: Bash(git:*), Bash(mkdir:*), Bash(cp:*), Bash(chmod:*), Bash(bash:*), Bash(make:*), Bash(cat:*), Bash(find:*), Bash(mise:*)
---

# Multiple Git Worktree Creator

TASK.mdファイルから複数のworktreeを一括作成し、並列開発環境を構築します。

## 概要

このSkillは以下のスクリプトを実行します：

```bash
bash .claude/skills/create-multiple-worktrees/scripts/create_multiple_worktrees.sh --base <branch> <task-files>
```

スクリプトにより以下の処理が行われます：

1. 引数で指定されたTASK.mdファイルを読み込み
2. ファイル名からfeature名を抽出（例: `user-auth.md` → `user-auth`）
3. 指定されたベースブランチから各featureの `.worktrees/<feature-name>/` にworktreeを作成
4. `feature/<feature-name>` ブランチを新規作成
5. TASK.mdを各worktreeにコピー
6. 環境変数ファイル（`.env*`, `.envrc`）を自動検出してコピー
7. 環境構築（mise install, make setup）を自動実行

## 使用方法

### 基本的な使い方

```bash
# ベースブランチ指定は必須
/create-multiple-worktrees --base main tasks/user-auth.md tasks/dashboard.md

# developブランチから作成
/create-multiple-worktrees --base develop tasks/*.md
```

### オプション

| オプション | 必須 | 説明 |
|------------|------|------|
| `--base <branch>` | **必須** | worktree作成元のブランチ（例: main, develop） |
| `--no-setup` | - | 環境構築（mise install, make setup）をスキップ |
| `--dry-run` | - | 実際には作成せず、何が作成されるか表示 |

## 実行結果

```
.worktrees/
├── user-auth/
│   ├── TASK.md          # ← tasks/user-auth.md からコピー
│   ├── .env*            # ← 自動検出してコピー
│   └── ...
├── dashboard/
│   ├── TASK.md          # ← tasks/dashboard.md からコピー
│   ├── .env*            # ← 自動検出してコピー
│   └── ...
└── api-v2/
    ├── TASK.md          # ← tasks/api-v2.md からコピー
    ├── .env*            # ← 自動検出してコピー
    └── ...
```

## 環境変数ファイルの自動検出

プロジェクトルートから3階層以内の `.env*` ファイルを自動検出してコピーします：

- `.env`, `.env.local`, `.env.dev`, `.env.test` など
- `apps/web/.env`, `apps/api/.env` など
- `.envrc`（direnv用）

除外されるディレクトリ: `node_modules`, `.venv`, `.worktrees`, `.git`

## 環境構築の自動実行

worktree作成後、以下の順で環境構築を試みます：

1. `mise.toml` / `.mise.toml` があれば → `mise install`
2. `Makefile` に `setup:` ターゲットがあれば → `make setup`

`--no-setup` オプションでスキップ可能です。

## 作業完了後

各worktreeで **create-pr** または **cleanup-worktree** スキルを使用：

```bash
cd .worktrees/<feature-name>
/create-pr              # PR作成
/cleanup-worktree       # worktree削除
```

### 手動でworktreeを削除する場合

```bash
# 個別に削除
git worktree remove .worktrees/<feature-name>

# 一括削除
for wt in .worktrees/*/; do git worktree remove "$wt"; done
```

## 詳細

詳細については [REFERENCE.md](REFERENCE.md) を参照してください。
