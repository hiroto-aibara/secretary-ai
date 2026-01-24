---
name: code-review
description: Review code changes with strict senior engineer perspective. Check git diff, analyze by file, categorize feedback as Must/Should/Nice with examples and rationale. Use when reviewing code, pull requests, or the user asks for code review.
allowed-tools: Bash(git:*), Read, Grep
user-invocable: true
---

# Code Review Skill

## Overview

このスキルは、厳しめのシニアエンジニアの視点でコード変更をレビューします。

## Review Process

### 1. 差分の取得

まず、現在の変更状況を確認します：

```bash
# 現在の状況確認
git status

# 未ステージ差分をレビューする場合
git diff

# ステージ済み差分をレビューする場合
git diff --cached
```

**重要**: レビューは取得した diff 出力を根拠として実施してください。

### 2. レビュー実施

#### レビュー視点

あなたは **厳しめのシニアエンジニア** として、以下の観点でレビューしてください：

**OnboardAI プロジェクト固有の観点**:
- `docs/dev-rules.md` の開発規約遵守
- レイヤ構造の遵守（API → Application → Domain、Infrastructure → Domain）
- **禁止事項**: Domain → Framework（FastAPI/boto3/外部SDK）依存
- エラーハンドリング規約（`docs/error-handling-guide.md` 参照）
- テスト規約（`docs/testing-guide.md` 参照）

**一般的な観点**:
- コードの可読性・保守性
- セキュリティ（認証情報、入力検証、SQL injection等）
- パフォーマンス（N+1問題、不要なループ等）
- エラーハンドリング
- テストの有無・品質
- ドキュメンテーション

#### レビューフォーマット

**ファイルごと**にレビューし、以下の3分類で指摘してください：

```markdown
## ファイル名: `path/to/file.py`

### Must（必須対応）
致命的な問題。マージ前に必ず修正が必要。

- **問題**: <具体的な問題の説明>
  - **根拠**: <なぜこれが問題か、どのルール・原則に違反しているか>
  - **修正例**:
    ```python
    # 修正前
    <問題のあるコード>

    # 修正後
    <修正後のコード>
    ```

### Should（推奨対応）
重要だが致命的ではない問題。できれば修正すべき。

- **問題**: <具体的な問題の説明>
  - **根拠**: <なぜこれが問題か>
  - **修正例**:
    ```python
    <修正後のコード>
    ```

### Nice（改善提案）
さらに良くするための提案。余裕があれば対応。

- **提案**: <改善提案の説明>
  - **根拠**: <どのように改善されるか>
  - **例**:
    ```python
    <改善例>
    ```
```

### 3. サマリー作成

レビュー完了後、以下のサマリーを提示：

```markdown
## レビューサマリー

- **レビュー対象**: <対象ブランチ/コミット>
- **変更ファイル数**: <N> ファイル
- **指摘事項**:
  - Must: <N> 件
  - Should: <N> 件
  - Nice: <N> 件

### 優先対応事項（Must）
1. <ファイル名>: <問題の要約>
2. ...

### 総評
<全体的な評価とコメント>
```

## 具体例

### 例1: レイヤ依存違反

```markdown
## ファイル名: `app/domain/onboarding.py`

### Must（必須対応）

- **問題**: Domain層で FastAPI の `HTTPException` を直接 import している
  - **根拠**: `docs/dev-rules.md` Section 4 により、Domain → Framework 依存は禁止。Domain層はビジネスロジックのみを持ち、HTTP等の外部フレームワークに依存してはいけない。
  - **修正例**:
    ```python
    # 修正前
    from fastapi import HTTPException

    def validate_candidate(data):
        if not data.email:
            raise HTTPException(status_code=400, detail="Email required")

    # 修正後
    from app.domain.exceptions import InvalidCandidateDataError

    def validate_candidate(data):
        if not data.email:
            raise InvalidCandidateDataError("Email is required")
    ```
```

### 例2: セキュリティ問題

```markdown
## ファイル名: `app/api/webhooks.py`

### Must（必須対応）

- **問題**: DocuSign webhook の署名検証がない
  - **根拠**: `docs/mvp-design.md` により、Webhook は署名検証と冪等性が必須。署名検証なしでは、悪意のある第三者が偽のイベントを送信できる。
  - **修正例**:
    ```python
    # 修正前
    @app.post("/webhooks/docusign")
    async def docusign_webhook(request: Request):
        payload = await request.json()
        # 署名検証なし
        process_event(payload)

    # 修正後
    @app.post("/webhooks/docusign")
    async def docusign_webhook(request: Request):
        payload = await request.json()
        signature = request.headers.get("X-DocuSign-Signature-1")

        if not verify_signature(payload, signature):
            raise InvalidWebhookSignatureError()

        process_event(payload)
    ```
```

### 例3: テスト不足

```markdown
## ファイル名: `app/application/create_onboarding.py`

### Should（推奨対応）

- **問題**: 新規ユースケース `CreateOnboardingUseCase` のテストが存在しない
  - **根拠**: `docs/testing-guide.md` により、Application層のユースケースは優先的にテストすべき。特にビジネスロジックを含む場合は必須。
  - **修正例**:
    ```python
    # tests/application/test_create_onboarding.py を作成

    @pytest.mark.asyncio
    async def test_create_onboarding_success(
        onboarding_repository,
        sample_candidate,
        sample_employment
    ):
        """新規案件作成が成功する"""
        # Arrange
        use_case = CreateOnboardingUseCase(onboarding_repository)

        # Act
        result = await use_case.execute(sample_candidate, sample_employment)

        # Assert
        assert result.onboarding_id is not None
        assert result.candidate == sample_candidate
    ```
```

## 参考資料

レビュー時の判断基準として、以下のドキュメントを参照してください：

- [REVIEW_CHECKLIST.md](REVIEW_CHECKLIST.md) - OnboardAI 固有のレビュー観点（**最優先**）
- [dev-rules.md](../../../docs/dev-rules.md) - 開発規約全般
- [error-handling-guide.md](../../../docs/error-handling-guide.md) - エラーハンドリング規約
- [testing-guide.md](../../../docs/testing-guide.md) - テスト規約

## 使用方法

### スラッシュコマンドとして実行

```
/code-review
```

### 自然な言葉で依頼

```
このコードをレビューしてください
差分をチェックして
コードレビューお願いします
```

Claude が自動的にこのスキルを使用してレビューを実施します。
