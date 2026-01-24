# PR Reviewer - Reference

## 段階的レビューの詳細手順

### Phase 1: 全体像把握

#### Step 1-1: PR基本情報取得

```bash
# MCPツール使用
mcp__github__get_pull_request(owner, repo, pullNumber)

# または gh CLI
gh pr view $PR_NUMBER --json number,title,body,author,baseRefName,headRefName,additions,deletions,changedFiles,mergeable,mergeableState
```

#### Step 1-2: 変更ファイル一覧取得

```bash
# MCPツール使用
mcp__github__get_pull_request_files(owner, repo, pullNumber)

# 結果が大きい場合（ファイルに保存された場合）
jq -r '.[0].text | fromjson | .[] | "\(.filename) (+\(.additions), -\(.deletions))"' <result_file>

# ファイル名のみ抽出
jq -r '.[0].text | fromjson | .[].filename' <result_file>
```

---

### Phase 2: 優先度判定

#### 優先度判定ルール（詳細）

```
優先度スコア計算:

1. セキュリティキーワード（+100）
   - auth, jwt, token, password, secret, credential, oauth, session
   - encrypt, decrypt, hash, sign, verify
   - permission, role, access, policy

2. パス優先度（+50〜+10）
   - api/, routes/, handlers/, endpoints/ → +50
   - domain/, use_cases/, services/, application/ → +40
   - infrastructure/, repositories/ → +30
   - components/, views/, pages/ → +20
   - tests/, __tests__/ → +10
   - *.md, *.json, *.yaml, *.yml → +5

3. 変更量（+0〜+20）
   - 追加行数 > 100 → +20
   - 追加行数 > 50 → +10
   - 追加行数 > 20 → +5

最終優先度 = セキュリティ + パス + 変更量
```

#### 分類の実行例

```markdown
入力ファイル一覧:
- apps/api/app/api/auth.py (+120, -10)
- apps/api/app/domain/user.py (+30, -5)
- apps/api/app/infrastructure/jwt_service.py (+50, -0)
- apps/web/src/views/LoginPage.tsx (+100, -20)
- tests/api/test_auth.py (+80, -0)
- README.md (+10, -2)

↓ 優先度計算後

1. [170] apps/api/app/infrastructure/jwt_service.py ← auth+jwt+infra+変更量
2. [160] apps/api/app/api/auth.py ← auth+api+変更量大
3. [50] apps/api/app/domain/user.py ← domain+変更量
4. [40] apps/web/src/views/LoginPage.tsx ← views+変更量大
5. [15] tests/api/test_auth.py ← tests+変更量
6. [5] README.md ← md
```

---

### Phase 3: 詳細レビュー

#### 個別diff取得

```bash
# 特定ファイルのdiffのみ
gh pr diff $PR_NUMBER -- path/to/file.py

# 複数ファイルを一度に（優先度高の上位N件）
gh pr diff $PR_NUMBER -- path/to/file1.py path/to/file2.py
```

#### 関連コード確認が必要なケース

```
1. インポート元の確認
   - 新規importがある場合、そのモジュールの実装を確認

2. インターフェース変更
   - Repository/Service のシグネチャ変更時、呼び出し元を確認

3. 型定義変更
   - 型の追加/変更時、使用箇所を確認

4. 設定変更
   - 環境変数追加時、.env.example を確認
```

#### 早期終了条件

```
以下の場合、残りのファイル確認をスキップ可能:

1. Critical問題が3つ以上検出
2. セキュリティ上の重大な脆弱性を検出
3. 明らかなマージ不可状態（テスト失敗、コンフリクト等）

→ Request Changes 確定として Phase 4 へ
```

---

## レビュー観点チェックリスト

### 1. コード品質

#### 可読性
- [ ] 変数名・関数名は意図を明確に表現しているか
- [ ] 複雑なロジックにはコメントがあるか
- [ ] 関数は単一責任になっているか
- [ ] ネストが深すぎないか（3段階以内）

#### 保守性
- [ ] DRY原則（重複コードがないか）
- [ ] マジックナンバーが定数化されているか
- [ ] 適切な抽象化レベルか
- [ ] 設定値がハードコードされていないか

#### 設計
- [ ] 責務の分離ができているか
- [ ] 依存関係が適切か（循環依存がないか）
- [ ] レイヤー構造が守られているか

### 2. バグ/ロジックエラー

#### 境界値・エッジケース
- [ ] 空配列・空文字列の処理
- [ ] None/null/undefined のチェック
- [ ] 0や負の数の処理
- [ ] 最大値・最小値の処理

#### 非同期処理
- [ ] async/await のエラーハンドリング
- [ ] 競合状態（Race Condition）の可能性
- [ ] タイムアウト処理
- [ ] リトライロジック

#### 型安全性
- [ ] 型アサーションの妥当性
- [ ] Any型の使用箇所
- [ ] Optional型の適切な処理

### 3. セキュリティ

#### インジェクション
- [ ] SQLインジェクション（プレースホルダ使用）
- [ ] XSS（出力エスケープ）
- [ ] コマンドインジェクション
- [ ] パストラバーサル

#### 認証・認可
- [ ] 認証チェックの漏れ
- [ ] 認可チェックの漏れ
- [ ] セッション管理の問題

#### 機密情報
- [ ] APIキー・トークンのハードコード
- [ ] ログへの機密情報出力
- [ ] エラーメッセージでの情報漏洩

### 4. パフォーマンス

#### データベース
- [ ] N+1クエリ問題
- [ ] 不要なカラム取得（SELECT *）
- [ ] インデックス活用
- [ ] 大量データの一括処理

#### アルゴリズム
- [ ] 計算量（O(n²)以上に注意）
- [ ] 不要なループ・再計算
- [ ] メモ化・キャッシュの活用

---

## レビューコメントテンプレート

### Critical（修正必須）

```markdown
🔴 **Critical: [問題タイトル]**

**場所**: `ファイル名:行番号`

**問題**:
[問題の詳細説明]

**リスク**:
[このまま放置した場合のリスク]

**修正案**:
[具体的な修正方法]
```

### Warning（要検討）

```markdown
🟡 **Warning: [問題タイトル]**

**場所**: `ファイル名:行番号`

**懸念点**:
[懸念の詳細]

**提案**:
[改善案]
```

### Suggestion（提案）

```markdown
🔵 **Suggestion: [提案タイトル]**

**場所**: `ファイル名:行番号`

**現状**: [現在の実装]
**提案**: [より良い実装案]
**理由**: [なぜこの変更が望ましいか]
```

---

## アクション選択基準

### Approve（承認）
- Critical問題がない
- Warning問題も軽微または許容範囲
- 全体的にコード品質が良い

### Request Changes（変更要求）
- Critical問題が1つ以上
- セキュリティ上の重大な懸念
- マージすると問題が発生する可能性が高い

### Comment（コメントのみ）
- Criticalはないが、Warning複数
- 判断に迷う点がある
- 追加の議論が必要

---

## トラブルシューティング

### Q: diffが大きすぎる

```bash
# ファイル単位で取得
gh pr diff $PR_NUMBER -- path/to/specific/file.py

# 統計のみ確認
gh pr view $PR_NUMBER --json additions,deletions,changedFiles
```

### Q: MCPツールの結果が大きすぎる

```bash
# 結果ファイルからjqで必要部分のみ抽出
jq -r '.[0].text | fromjson | .[].filename' <result_file>

# 追加行数でフィルタ（大きな変更のみ）
jq -r '.[0].text | fromjson | .[] | select(.additions > 50) | .filename' <result_file>
```

### Q: 特定ファイルの完全な内容を確認したい

```bash
# ローカルでPRブランチをfetch
git fetch origin pull/$PR_NUMBER/head:pr-$PR_NUMBER

# 特定ファイルを確認
git show pr-$PR_NUMBER:path/to/file.py

# またはReadツールでHEADブランチのファイルを読む
```

---

## ベストプラクティス

### 1. 効率的なレビュー
- 全体像を把握してから詳細に入る
- 優先度の高いファイルから確認
- 早期終了条件を活用

### 2. 建設的なフィードバック
- 問題点だけでなく良い点も指摘
- 「なぜ」を説明する
- 具体的な改善案を提示

### 3. 適切な粒度
- 重要な問題に集中
- 些細なスタイル指摘は控えめに
- 別PRで対応すべきものは分ける

---

## 関連ドキュメント

- [SKILL.md](SKILL.md) - 基本的な使い方
- [GitHub CLI Documentation](https://cli.github.com/manual/)
- [GitHub Pull Request Reviews](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests)
