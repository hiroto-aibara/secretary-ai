---
name: create-pr
description: Creates a PR from the current branch without deleting the worktree. Use this when you want to keep the worktree for local modifications.
allowed-tools: Bash(git:*), Bash(gh:*), Bash(cd:*), Bash(pwd:*), Read
---

# Create PR

現在のブランチからPRを作成します。

## 概要

このスキルは以下のワークフローを実行します：

```
1. 環境検証（bash）
   - worktree/ブランチ検出
   - 未コミット変更チェック
   - 未プッシュコミット → 自動push
        ↓
2. 情報収集・構造化（bash）【Token節約】
   - レイヤー別変更ファイル集計
   - コミットログ要約
   - TASK.md抽出（概要・受け入れ基準）
   → 構造化情報を出力
        ↓
3. PR本文生成（Claude Code）
   - 構造化情報のみを入力として受け取る
   - テンプレートに沿って本文生成
        ↓
4. PR作成（bash）
   - gh pr create 実行
   - PR URL返却
```

## 使用方法

### 基本的な使い方

```bash
# worktreeディレクトリ内で実行
cd .worktrees/<feature-name>
/create-pr
```

### オプション

```bash
# タイトルと本文を事前指定
/create-pr --title "feat: Add new feature" --body "詳細な説明..."

# ドラフトPRとして作成
/create-pr --draft

# ベースブランチを指定
/create-pr --base develop
```

---

## PR本文の書き方ガイド

### 必須セクション

#### 1. Summary（1-3行）

変更の目的と効果を簡潔に記述。

**良い例:**
> オンボーディング案件の一覧・詳細取得APIを実装し、フロントエンドをAPI経由のデータ取得に移行。
> これによりページリロード後もデータが保持されるようになる。

**悪い例:**
> APIを追加しました。

#### 2. 変更内容（レイヤー別）

変更をレイヤー別に整理し、各変更の具体的な内容を記述。

**良い例:**
```markdown
### Backend
- **Domain層**: `OnboardingRepository` に `list_all()`, `count_all()` 追加
- **Infrastructure層**: DynamoDB Scan実装（将来的にGSI移行予定）
- **Application層**: `GetOnboardingUseCase`, `ListOnboardingsUseCase` 実装
- **API層**: `GET /onboardings`, `GET /onboardings/{id}` エンドポイント追加

### Frontend
- APIクライアント関数追加（`getOnboarding`, `listOnboardings`）
- OnboardingListPage/DetailPage をAPI経由に移行
- Zustandストアからの依存を削除
```

**悪い例:**
```markdown
- ファイルを追加
- 修正
- テスト追加
```

#### 3. Test plan

検証方法をチェックリスト形式で記述。

**良い例:**
```markdown
- [ ] `pytest tests/application/` - UseCase単体テスト
- [ ] `pytest tests/api/` - API統合テスト
- [ ] `ruff check && mypy` - Lintエラーなし
- [ ] フロントエンドで一覧・詳細画面が表示される
- [ ] ページリロード後もデータが保持される
```

### TASK.mdとの連携

worktree内に `TASK.md` がある場合、スクリプトが以下を自動抽出：
- **概要** → Summary の素材
- **要件** → 変更内容の構造
- **受け入れ基準** → Test plan のベース

---

## 前提条件

- worktreeディレクトリ内で実行すること
- すべての変更がコミット済みであること
- `gh` CLI がインストール・認証済みであること

## 実行例

```bash
$ cd .worktrees/user-auth
$ /create-pr

[INFO] Current branch: feature/user-auth
[INFO] Base branch: main
[STEP] Analyzing changes...

### 変更ファイル分析

**Backend:**
- Domain層: 1ファイル
- Application層: 2ファイル
- Infrastructure層: 1ファイル
- API層: 1ファイル

**Frontend:**
- 4ファイル変更

**テスト:**
- 6ファイル追加

### TASK.md情報
- Summary候補: ユーザー認証機能を実装...
- Test plan候補: 全テストがパス, ESLintエラーなし...

[STEP] Creating pull request...
[INFO] PR created: https://github.com/user/repo/pull/123
```

## 関連スキル

- `cleanup-worktree`: worktree削除
- `review-pr`: PRレビュー
- `create-worktree`: worktree作成

## 詳細

詳細については [REFERENCE.md](REFERENCE.md) を参照してください。
