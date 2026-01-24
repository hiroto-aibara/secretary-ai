---
name: review-pr
description: Review a GitHub PR and post comments. Analyzes code quality, bugs, security, and performance. User chooses final action (Approve/Request Changes/Comment).
allowed-tools: Bash(gh:*), Bash(git:*), Bash(jq:*), Read, Grep, Glob, mcp__github__get_pull_request, mcp__github__get_pull_request_files, mcp__github__get_pull_request_diff, mcp__github__create_and_submit_pull_request_review
---

# PR Reviewer

GitHub PRを段階的にレビューし、コメントを投稿します。

## 概要

このスキルは**段階的アプローチ**でPRをレビューします：

```
Phase 1: 全体像把握（低Token消費）
   - PR基本情報（タイトル、説明、統計）
   - 変更ファイル一覧（パスと行数のみ）
   - レイヤー別分類・優先度判定
        ↓
Phase 2: 優先度判定
   - セキュリティ関連ファイルを特定
   - 重要度でソート
   - レビュー順序を決定
        ↓
Phase 3: 詳細レビュー（ファイル単位で段階的）
   - 優先度高のファイルから個別にdiff取得
   - 必要に応じて関連コードを読み込み
   - 問題点を特定
        ↓
Phase 4: レビュー結果整理・提示
        ↓
Phase 5: ユーザー確認・投稿
```

## 使用方法

### 基本的な使い方

```bash
# PR番号で指定
/review-pr 123

# PR URLで指定
/review-pr https://github.com/owner/repo/pull/123
```

### owner/repoを明示的に指定

```bash
# 別リポジトリのPRをレビュー
/review-pr owner/repo#123
```

---

## Phase 1: 全体像把握

### 取得する情報

1. **PR基本情報**: `mcp__github__get_pull_request`
   - タイトル、説明、作成者
   - ベースブランチ、ヘッドブランチ
   - マージ可能状態

2. **変更ファイル一覧**: `mcp__github__get_pull_request_files`
   - 結果が大きい場合は `jq` でファイル名・統計のみ抽出:
   ```bash
   jq -r '.[0].text | fromjson | .[] | "\(.filename) (+\(.additions), -\(.deletions))"' <result_file>
   ```

### 出力フォーマット

```markdown
### PR概要
| 項目 | 内容 |
|------|------|
| タイトル | feat: Add authentication |
| 作成者 | username |
| ベース | main ← feature/auth |
| 状態 | mergeable: clean |
| ファイル数 | 15 |
| 変更 | +500 / -100 |

### 変更ファイル一覧
- app/api/auth.py (+120, -10)
- app/domain/user.py (+30, -5)
- tests/test_auth.py (+80, -0)
...
```

---

## Phase 2: 優先度判定

### 優先度ルール

| 優先度 | 条件 | パターン例 |
|--------|------|------------|
| 🔴 最高 | セキュリティ関連 | auth, jwt, password, secret, credential, token |
| 🟠 高 | API/エンドポイント | api/, routes/, handlers/, endpoints/ |
| 🟡 中 | ドメイン/ビジネスロジック | domain/, use_cases/, services/, application/ |
| 🟢 低 | インフラ/リポジトリ | infrastructure/, repositories/ |
| 🔵 通常 | フロントエンド | components/, views/, pages/ |
| ⚪ 最低 | テスト/設定 | tests/, *.test.*, *.md, *.json, *.yaml |

### 分類出力例

```markdown
### レビュー優先順位

**🔴 セキュリティ関連（最優先）**
1. app/infrastructure/jwt_service.py (+50, -0)
2. app/api/auth.py (+120, -10)

**🟠 API層**
3. app/api/onboardings.py (+80, -20)

**🟡 Application/Domain層**
4. app/use_cases/login.py (+60, -0)
5. app/domain/user.py (+30, -5)

**⚪ テスト/その他**
6. tests/test_auth.py (+80, -0)
7. README.md (+10, -2)
```

---

## Phase 3: 詳細レビュー

### 進め方

優先度順にファイルを確認：

```
for each file in priority_order:
    1. gh pr diff <PR> -- <file> で個別diff取得
    2. 4つの観点でチェック
    3. 必要なら関連ファイルを Read で確認
    4. 問題点をリストに追加

    # 早期終了条件
    if Critical問題が3つ以上:
        → これ以上の詳細確認は不要
        → Request Changes確定として Phase 4 へ
```

### 個別diff取得コマンド

```bash
# 特定ファイルのdiffのみ取得
gh pr diff <PR_NUMBER> -- path/to/file.py
```

### レビュー観点（4つ）

1. **コード品質**: 可読性、保守性、設計パターン
2. **バグ/ロジックエラー**: 境界値、null処理、非同期エラー
3. **セキュリティ**: インジェクション、認証認可、機密情報
4. **パフォーマンス**: N+1、不要ループ、メモリリーク

---

## Phase 4: レビュー結果整理

### 結果フォーマット

```markdown
## PR Review: #123 - PRタイトル

### 📋 概要
| 項目 | 内容 |
|------|------|
| 変更ファイル数 | 15 files |
| 追加/削除 | +500 / -100 |
| レビュー対象 | 8 files（優先度高のみ詳細確認） |

**変更の概要**: [1-2文で変更内容を要約]

---

### 🔍 レビュー結果

#### 🔴 Critical（修正必須）
なし / または問題リスト

#### 🟡 Warning（要検討）
なし / または問題リスト

#### 🔵 Suggestion（提案）
なし / または提案リスト

---

### ✅ 良い点
- [評価点]

---

### 📊 総合評価
**推奨アクション**: Approve / Request Changes / Comment
**理由**: [1-2文]
```

---

## Phase 5: ユーザー確認・投稿

### アクション選択基準

| アクション | 条件 |
|------------|------|
| **Approve** | Critical なし、Warning 軽微 |
| **Request Changes** | Critical 1つ以上、またはセキュリティ懸念 |
| **Comment** | 判断に迷う、追加議論が必要 |

### Request Changes 時のオプション

1. **通常のRequest Changes**: レビューコメントのみ投稿
2. **@claude 付き**: GitHub Actionで自動修正させる場合

---

## Token消費の目安

| PR規模 | 従来 | 改善後 |
|--------|------|--------|
| 小（〜10ファイル） | 5,000 | 1,500 |
| 中（〜30ファイル） | 15,000 | 3,000 |
| 大（50ファイル超） | 30,000+ | 5,000 |

---

## 前提条件

- `gh` CLI がインストール・認証済み
- PRの読み取り権限があること

## 関連スキル

- `create-pr`: PR作成
- `cleanup-worktree`: worktree削除

## 詳細

詳細については [REFERENCE.md](REFERENCE.md) を参照してください。
