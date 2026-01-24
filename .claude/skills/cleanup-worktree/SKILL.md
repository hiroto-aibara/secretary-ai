---
name: cleanup-worktree
description: Removes the current worktree and returns to the main branch. Use after PR is merged.
allowed-tools: Bash(git:*), Bash(cd:*), Bash(pwd:*)
---

# Cleanup Worktree

現在のworktreeを削除し、mainブランチに戻ります。

## 概要

このスキルは以下を実行します：

1. 現在のworktreeディレクトリを検出
2. リポジトリルートに移動
3. worktreeを削除
4. mainブランチにチェックアウト

## 使用方法

### 基本的な使い方

```bash
# worktreeディレクトリ内で実行
cd .worktrees/<feature-name>
/cleanup-worktree
```

### オプション

```bash
# 復帰先ブランチを指定
/cleanup-worktree --branch develop

# 強制削除（未コミット変更があっても削除）
/cleanup-worktree --force

# ローカルブランチも削除
/cleanup-worktree --delete-branch
```

## 実行例

```bash
$ cd .worktrees/user-auth
$ /cleanup-worktree

[STEP] Removing worktree...
[INFO] Worktree removed: .worktrees/user-auth
[STEP] Returning to main branch...
[INFO] Now on branch: main
```

## 注意事項

- ローカルブランチは保持されます（--delete-branch を指定しない限り）
- リモートブランチは保持されます
- 未コミット変更がある場合は警告が表示されます

## 関連スキル

- `create-pr`: PR作成
- `create-worktree`: worktree作成

## 詳細

詳細については [REFERENCE.md](REFERENCE.md) を参照してください。
