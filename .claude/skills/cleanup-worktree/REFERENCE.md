# Cleanup Worktree - Reference

## コマンドライン引数

| 引数 | 説明 | デフォルト |
|------|------|------------|
| `--branch <branch>` | 復帰先ブランチ | main |
| `--force` | 未コミット変更があっても強制削除 | false |
| `--delete-branch` | ローカルブランチも削除 | false |
| `--help` | ヘルプメッセージを表示 | - |

## 内部処理フロー

### 1. 環境検証フェーズ
```bash
1. worktreeディレクトリ内かチェック
   - git rev-parse --git-dir で .git/worktrees/<name> を確認
2. 現在のブランチ名取得
   - git branch --show-current
3. リポジトリルート取得
   - git worktree list で親リポジトリを特定
```

### 2. 変更チェックフェーズ
```bash
1. 未コミット変更チェック
   - git status --porcelain
   - 出力が空でなければ警告（--force で無視可能）
```

### 3. クリーンアップフェーズ
```bash
1. リポジトリルートに移動
   - cd $REPO_ROOT
2. worktree削除
   - git worktree remove <worktree-path>
3. 復帰先ブランチにチェックアウト
   - git checkout $BRANCH
4. ローカルブランチ削除（オプション）
   - git branch -d $FEATURE_BRANCH
```

## 使用例

### ケース1: 基本的な使い方
```bash
cd .worktrees/feature-x
bash ../../.claude/skills/cleanup-worktree/scripts/cleanup_worktree.sh
# → worktree削除 → mainに戻る
```

### ケース2: developブランチに戻る
```bash
bash cleanup_worktree.sh --branch develop
```

### ケース3: 強制削除
```bash
# 未コミット変更があっても削除
bash cleanup_worktree.sh --force
```

### ケース4: ローカルブランチも削除
```bash
# PRマージ後、ブランチが不要な場合
bash cleanup_worktree.sh --delete-branch
```

## トラブルシューティング

### Q: worktree削除に失敗する
```bash
# 手動でworktreeを削除
cd /path/to/repo-root
git worktree remove .worktrees/<feature-name>

# 強制削除
git worktree remove --force .worktrees/<feature-name>
```

### Q: worktree外で実行してしまった
```bash
# worktreeディレクトリに移動
cd .worktrees/<feature-name>

# 再実行
bash ../../.claude/skills/cleanup-worktree/scripts/cleanup_worktree.sh
```

### Q: ローカルブランチを後から削除したい
```bash
# マージ済みブランチを削除
git branch -d feature/<name>

# 強制削除（マージされていなくても）
git branch -D feature/<name>
```

### Q: 複数のworktreeを一括削除したい
```bash
# 一括削除
for wt in .worktrees/*/; do git worktree remove "$wt"; done
```

## 関連ドキュメント

- [SKILL.md](SKILL.md) - 基本的な使い方
- [create-pr スキル](../create-pr/SKILL.md) - PR作成
- [Git Worktree Documentation](https://git-scm.com/docs/git-worktree) - git worktree公式ドキュメント
