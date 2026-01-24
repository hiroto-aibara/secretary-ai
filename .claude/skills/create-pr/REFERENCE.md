# Create PR - Reference

## コマンドライン引数

| 引数 | 説明 | デフォルト |
|------|------|------------|
| `--title <title>` | PRタイトル | インタラクティブ入力 |
| `--body <body>` | PR本文 | インタラクティブ入力 |
| `--base <branch>` | ベースブランチ | main |
| `--draft` | ドラフトPRとして作成 | false |
| `--force` | 未コミット変更があっても続行 | false（非推奨） |
| `--help` | ヘルプメッセージを表示 | - |

## 内部処理フロー

### 1. 環境検証フェーズ
```bash
1. gitリポジトリ内かチェック
   - git rev-parse --git-dir
2. worktreeディレクトリ内かチェック
   - .git/worktrees/<name> を確認
3. 現在のブランチ名取得
   - git branch --show-current
4. gh CLI のインストール確認
```

### 2. 変更チェックフェーズ
```bash
1. 未コミット変更チェック
   - git status --porcelain
   - 出力が空でなければエラー
2. push状態チェック
   - git rev-list origin/<branch>..HEAD
   - unpushedコミットがあれば自動push
```

### 3. 情報収集・構造化フェーズ【Token節約】
```bash
1. レイヤー別ファイル集計
   - git diff --name-only $BASE_BRANCH...HEAD
   - app/domain/, app/application/, app/infrastructure/, app/api/ を分類
   - apps/web/ をフロントエンドとして集計
   - tests/ をテストとして集計

2. コミットログ要約
   - git log --oneline $BASE_BRANCH..HEAD

3. TASK.md情報抽出（存在する場合）
   - 概要セクション → Summary候補
   - 受け入れ基準セクション → Test plan候補

4. 統計情報
   - git diff --stat $BASE_BRANCH...HEAD
```

### 4. PR作成フェーズ
```bash
1. gh pr create 実行
   - --title / --body 指定がなければインタラクティブ
   - --base でベースブランチ指定（デフォルト: main）
2. PR URL取得
   - gh pr view --json url -q .url
3. 成功確認
   - 終了コード 0 を確認
```

---

## 構造化情報フォーマット

スクリプトが出力する構造化情報の例：

```
### 変更ファイル分析

**Backend:**
- Domain層: 1ファイル (repositories.py)
- Application層: 2ファイル (get_onboarding.py, list_onboardings.py)
- Infrastructure層: 1ファイル (onboarding_repository.py)
- API層: 1ファイル (onboardings.py)

**Frontend:**
- 4ファイル変更 (api/, views/)

**テスト:**
- 6ファイル追加

**統計:**
- 25ファイル変更, +1963行, -328行

**コミット履歴:**
- feat: Add GetOnboardingUseCase
- feat: Add ListOnboardingsUseCase
- feat: Add API endpoints

### TASK.md情報

**Summary候補:**
オンボーディング案件の一覧取得・詳細取得APIを実装し、フロントエンドと統合する。

**Test plan候補:**
- [ ] GET /onboardings で案件一覧が取得できる
- [ ] GET /onboardings/{id} で案件詳細が取得できる
- [ ] 存在しないIDで404エラーが返る
- [ ] 全テストがパス
```

---

## PR本文テンプレート

Claude Codeが生成するPR本文のテンプレート：

```markdown
## Summary
{task.summaryまたはcommitsから1-3行で生成}

## 変更内容

### Backend
- **Domain層**: {具体的な変更内容}
- **Application層**: {具体的な変更内容}
- **Infrastructure層**: {具体的な変更内容}
- **API層**: {具体的な変更内容}

### Frontend
- {具体的な変更内容}

### Tests
- {追加したテストの概要}

## Test plan
{task.acceptanceをチェックリスト化、なければ標準項目}

- [ ] 単体テストがパス
- [ ] 統合テストがパス
- [ ] Lint/型チェックがパス
- [ ] 動作確認完了

---

🤖 Generated with [Claude Code](https://claude.ai/code)
```

---

## 使用例

### ケース1: 基本的な使い方
```bash
cd .worktrees/feature-x
/create-pr
# → 情報収集 → PR本文生成 → PR作成
```

### ケース2: タイトル・本文を事前指定
```bash
/create-pr \
  --title "feat(frontend): Add user authentication" \
  --body "Implements Firebase Auth integration"
# → 情報収集スキップ → 指定内容でPR作成
```

### ケース3: ドラフトPR作成
```bash
/create-pr --draft
```

### ケース4: カスタムベースブランチ
```bash
# develop ブランチにPR作成
/create-pr --base develop
```

---

## トラブルシューティング

### Q: 未コミットの変更があると言われる
```bash
# 変更を確認
git status

# コミットする
git add .
git commit -m "commit message"

# または stash
git stash
```

### Q: gh CLI の認証が必要と言われる
```bash
# GitHub CLIの認証
gh auth login

# 認証状態確認
gh auth status
```

### Q: リモートブランチにpush済みか確認したい
```bash
# ローカルとリモートの差分確認
git rev-list origin/<branch>..HEAD --count

# 0なら差分なし
```

### Q: TASK.mdが読み込まれない
- `TASK.md` はworktreeルートに配置
- ファイル名は大文字小文字を正確に

---

## 関連ドキュメント

- [SKILL.md](SKILL.md) - 基本的な使い方
- [cleanup-worktree スキル](../cleanup-worktree/SKILL.md) - worktree削除
- [GitHub CLI Documentation](https://cli.github.com/manual/) - gh コマンドリファレンス
