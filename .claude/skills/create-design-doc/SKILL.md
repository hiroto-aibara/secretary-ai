---
name: create-design-doc
description: Creates a Design Doc with layer structure, domain interfaces, infrastructure, application use cases, API specs, and test strategy. Use after Feature Brief approval.
allowed-tools: Write, Read, Bash(mkdir:*), AskUserQuestion, Glob, Grep
---

# Design Doc Creator

Feature Brief を元に設計書（Design Doc）を生成します。

## 概要

このSkillは以下を実行します：

1. Feature Brief とコードベースを参照して設計情報を収集
2. レイヤ構成・インターフェース・実装仕様を含む Design Doc を生成
3. `docs/<feature-name>-design.md` として保存

## 使用方法

```bash
/create-design-doc <feature-name>
```

### 例

```bash
/create-design-doc user-auth
# → docs/user-auth-design.md が生成される
```

## Design Doc の位置づけ

```
Feature Brief（なぜ・何を）  ← /create-feature-brief で作成
    ↓
Design Doc（どうやって）     ← このスキルで作成
    ↓
Task File（実装指示）        ← /create-task で作成
```

## 実行手順

1. **情報収集**
   - `docs/<feature-name>-brief.md` を読む（存在すれば）
   - 既存コードベースのパターンを確認
   - 不明点があればユーザーに質問

2. **Design Doc 生成**
   - [TEMPLATE.md](TEMPLATE.md) のフォーマットに従って生成
   - 各レイヤの具体的なインターフェース・クラス・メソッドを定義
   - テスト戦略と実装順序を含める
   - セクション番号を付与（タスクファイルからの参照用）

3. **ユーザーレビュー**
   - 生成したドキュメントを表示
   - 修正リクエストがあれば反映

## 生成ファイル

```
docs/
├── user-auth-design.md
├── notification-design.md
└── ...
```

## セクション構成

Design Doc は以下のセクションで構成されます（詳細は TEMPLATE.md 参照）：

1. 概要
2. レイヤ構成
3. Domain 層（インターフェース、型、例外、Config）
4. Infrastructure 層（認証、クライアント、実装）
5. Application 層（UseCase、Orchestrator）
6. API エンドポイント
7. テスト戦略
8. 実装順序と依存関係
9. セキュリティ考慮事項
10. 環境変数
11. 依存パッケージ

## カスタマイズ

プロジェクトに応じてセクションを追加・削除してください。例：
- MCP ツールがある場合 → MCPセクション追加
- Bootstrap DI がある場合 → DI配線セクション追加
- 監査ログがある場合 → 監査ログセクション追加

## テンプレート

[TEMPLATE.md](TEMPLATE.md) にフォーマットがあります。

## 関連スキル

- `/create-feature-brief` - Feature Brief 作成（Design Doc の前段階）
- `/create-task` - タスクファイル作成（Design Doc のセクションを参照）
