# Task Creator - Reference

## TASK.mdの役割

TASK.mdは各Claude Codeセッションに対して以下を伝達します：

1. **タスクの目的** - 何を達成すべきか
2. **要件** - 具体的に実装すべき機能
3. **品質基準** - 守るべき基準（テスト、パフォーマンス、セキュリティ）
4. **受け入れ基準** - 完了条件

## 命名規則

| タスク種別 | ファイル名例 |
|-----------|-------------|
| 新機能 | `feature-user-auth.md` |
| バグ修正 | `fix-login-error.md` |
| リファクタリング | `refactor-api-client.md` |
| パフォーマンス改善 | `perf-query-optimization.md` |

## 良いTASK.mdの例

```markdown
# Task: user-auth

## 概要
JWTを使用したユーザー認証機能を実装する。
セキュアなログイン・ログアウト・トークン管理を提供。

## 要件
- [ ] POST /api/auth/login - メール/パスワードでログイン
- [ ] POST /api/auth/logout - ログアウト（トークン無効化）
- [ ] POST /api/auth/refresh - アクセストークン更新
- [ ] GET /api/auth/me - 現在のユーザー情報取得
- [ ] パスワードハッシュ化（bcrypt）
- [ ] JWTトークン生成・検証

## 技術仕様
- 使用技術: Express.js, jsonwebtoken, bcrypt
- 関連ファイル: src/routes/auth.ts, src/middleware/auth.ts
- APIエンドポイント: /api/auth/*

## 受け入れ基準
- [ ] 全APIエンドポイントのユニットテスト
- [ ] 認証フローのE2Eテスト
- [ ] ESLint/Prettier エラーなし
- [ ] セキュリティレビュー完了

## 品質基準
- テストカバレッジ: 90%以上（認証は重要機能のため）
- パフォーマンス: ログインAPI 100ms以内
- セキュリティ:
  - パスワードは平文で保存しない
  - トークンは適切な有効期限を設定
  - レート制限を実装

## 参考リソース
- Issue #42: ユーザー認証の実装
- RFC 7519: JSON Web Token (JWT)
- OWASP Authentication Cheat Sheet
```

## Claude Codeへの指示例

TASK.md作成時にClaude Codeに伝える内容の例：

### 良い例
```
JWTを使った認証機能を実装したい。
- ログイン、ログアウト、トークンリフレッシュが必要
- パスワードはbcryptでハッシュ化
- アクセストークンは15分、リフレッシュトークンは7日で有効期限切れ
- 既存のExpressアプリに追加する形で
```

### 悪い例
```
認証機能作って
```

## 複数タスクの分割指針

大きなタスクは適切に分割してください：

### 分割すべき場合
- 異なる開発者が並列で作業できる
- 独立してテスト・デプロイできる
- 1つのPRが大きくなりすぎる（目安: 500行以上）

### 分割の例
```
❌ 悪い例: user-management.md（巨大すぎる）

✅ 良い例:
  - user-auth.md（認証）
  - user-profile.md（プロフィール管理）
  - user-settings.md（設定）
```

## ディレクトリ構成

```
project/
├── tasks/                    # TASK.mdファイル格納
│   ├── user-auth.md
│   ├── dashboard.md
│   └── api-v2.md
├── .worktrees/               # worktree（create-multiple-worktreesで作成）
│   ├── user-auth/
│   │   └── TASK.md          # tasks/user-auth.md のコピー
│   ├── dashboard/
│   │   └── TASK.md
│   └── api-v2/
│       └── TASK.md
└── .claude/
    └── skills/
        ├── create-task/
        └── create-multiple-worktrees/
```

## 関連スキル

- [create-multiple-worktrees](../create-multiple-worktrees/SKILL.md) - TASK.mdからworktree作成
- [pr-and-cleanup](../pr-and-cleanup/SKILL.md) - PR作成とworktree削除
