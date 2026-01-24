# Multiple Worktree Creator - Reference

## 使用例

### 複数のTASK.mdからworktree作成

```bash
bash .claude/skills/create-multiple-worktrees/scripts/create_multiple_worktrees.sh \
  --base main \
  tasks/user-auth.md \
  tasks/dashboard.md \
  tasks/api-v2.md

# 出力例:
# [INFO] Creating 3 worktrees from TASK.md files...
# [INFO] Base branch: main
# [STEP] Creating worktree: user-auth
# [INFO] TASK.md copied to .worktrees/user-auth
# [INFO] Detecting environment files...
# [INFO]   Copied: .env
# [INFO]   Copied: apps/web/.env.local
# [INFO] Copied 2 environment file(s)
# [INFO] Running environment setup...
# [INFO]   Running mise install...
# [✓] user-auth created (.worktrees/user-auth)
# ...
```

### developブランチから作成

```bash
bash .claude/skills/create-multiple-worktrees/scripts/create_multiple_worktrees.sh \
  --base develop \
  tasks/*.md
```

### ドライラン（事前確認）

```bash
bash .claude/skills/create-multiple-worktrees/scripts/create_multiple_worktrees.sh \
  --base main \
  --dry-run \
  tasks/*.md

# 出力例:
# [DRY-RUN] Would create the following worktrees:
# [DRY-RUN] Base branch: main
#
# [DRY-RUN] Detected environment files:
#     .env
#     apps/web/.env.local
#     apps/api/.env
#
#   - .worktrees/user-auth/
#     ├── TASK.md (from tasks/user-auth.md)
#     ├── .env* (auto-detected)
#     └── (branch: feature/user-auth from main)
```

### 環境構築をスキップ

```bash
bash .claude/skills/create-multiple-worktrees/scripts/create_multiple_worktrees.sh \
  --base main \
  --no-setup \
  tasks/quick-fix.md
```

## コマンドライン引数

| 引数 | 必須 | 説明 | デフォルト |
|------|------|------|------------|
| `--base <branch>` | **必須** | worktree作成元のブランチ | - |
| `--no-setup` | - | 環境構築をスキップ | false |
| `--dry-run` | - | 事前確認モード | false |
| `-h, --help` | - | ヘルプ表示 | - |

## TASK.mdファイル名とworktreeの対応

ファイル名がそのままfeature名・ブランチ名になります：

```
tasks/user-auth.md     → .worktrees/user-auth/  (branch: feature/user-auth)
tasks/fix-login.md     → .worktrees/fix-login/  (branch: feature/fix-login)
tasks/refactor-api.md  → .worktrees/refactor-api/ (branch: feature/refactor-api)
```

## 環境変数ファイルの自動検出

以下のルールで `.env*` ファイルを検出・コピーします：

### 検出対象
- プロジェクトルートから3階層以内
- `.env*` にマッチするファイル（`.env`, `.env.local`, `.env.dev` など）
- `.envrc`（direnv用）

### 除外対象
- `node_modules/` 以下
- `.venv/` 以下
- `.worktrees/` 以下
- `.git/` 以下

### ポートのランダム化
ルートの `.env` ファイルは、以下の変数がランダムなポートに置換されます：
- `FRONTEND_PORT`
- `BACKEND_PORT`
- `AGENT_PORT`
- `PORT`

## 環境構築の自動実行

以下の順序で環境構築を試みます：

1. **mise**: `mise.toml` または `.mise.toml` が存在する場合
   ```bash
   mise install
   ```

2. **Makefile**: `setup:` ターゲットが存在する場合
   ```bash
   make setup
   ```

## トラブルシューティング

### Q: --base を指定し忘れた

```
[ERROR] --base option is required. Specify the base branch (e.g., --base main)
```

`--base main` または `--base develop` を追加してください。

### Q: ベースブランチが存在しない

```
[ERROR] Base branch 'feature-x' does not exist
```

`git branch -a` で利用可能なブランチを確認してください。

### Q: TASK.mdファイルが見つからない

```bash
ls -la tasks/
```

### Q: worktree作成に失敗した

```bash
# 状態確認
git worktree list

# 既存のworktreeを削除して再作成
git worktree remove .worktrees/<feature-name>
```

### Q: ポートが競合する

各worktreeの `.env` にはランダムなポートが割り当てられます。
手動で変更する場合：

```bash
vim .worktrees/<feature-name>/.env
```

### Q: worktreeを削除したい

```bash
# 個別に削除
git worktree remove .worktrees/<feature-name>

# 一括削除
for wt in .worktrees/*/; do git worktree remove --force "$wt"; done
git worktree prune
```

### Q: mise install が失敗する

mise がインストールされているか確認：
```bash
which mise
mise --version
```

インストールされていない場合は [mise公式サイト](https://mise.jdx.dev/) を参照。

## 関連ドキュメント

- [Git Worktree公式ドキュメント](https://git-scm.com/docs/git-worktree)
- [mise公式ドキュメント](https://mise.jdx.dev/)
