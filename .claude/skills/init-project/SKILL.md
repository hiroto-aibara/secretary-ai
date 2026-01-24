---
name: init-project
description: Initialize a new project with git, GitHub repo, tooling foundation, and documentation templates. Run this first before init-go-backend or init-react-frontend.
allowed-tools: Bash, Write, Read, Edit
---

# Init Project

プロジェクトの基盤をセットアップします（技術スタック非依存）。

## 概要

このスキルは以下を実行します：

```
1. 情報収集（対話）
   - プロジェクト名
   - GitHub ユーザー/Organization
   - プロジェクト概要（1-2文）
        ↓
2. Git初期化
   - git init
   - .gitignore 作成
        ↓
3. プロジェクト基盤ファイル作成
   - mise.toml（スケルトン）
   - .editorconfig
   - CLAUDE.md（スケルトン）
        ↓
4. ドキュメントテンプレート作成
   - docs/architecture.md（テンプレート）
   - docs/api.md（テンプレート）
   - docs/conventions.md（テンプレート）
        ↓
5. ツーリング基盤
   - package.json（husky + lint-staged）
   - .husky/pre-commit
   - .github/dependabot.yml
        ↓
6. GitHub リポジトリ作成 + プッシュ
        ↓
7. Skill-source submodule 追加
   - .claude/skill-source（submodule）
   - .claude/skills/ に必要なスキルをコピー
```

## 使用方法

```bash
/init-project
```

## 作成されるファイル

```
.
├── .editorconfig
├── .gitignore
├── .github/
│   └── dependabot.yml
├── .husky/
│   └── pre-commit
├── .claude/
│   ├── skill-source/          (submodule)
│   └── skills/                (コピー)
├── docs/
│   ├── architecture.md        (テンプレート)
│   ├── api.md                 (テンプレート)
│   └── conventions.md         (テンプレート)
├── mise.toml
├── CLAUDE.md
└── package.json
```

## 各ファイルの内容

### .editorconfig

```ini
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
```

### mise.toml（スケルトン）

```toml
[tools]
# 技術スタックに応じて init-backend / init-frontend が追加

[tasks.clean]
description = "Clean build artifacts"
run = "rm -rf bin/ dist/"
```

### CLAUDE.md（スケルトン）

プロジェクト概要とドキュメントへのリンクを含む。
技術スタック固有の情報は init-backend / init-frontend が追記する。

### docs/ テンプレート

各ファイルにはセクション見出しのみを含むテンプレートを生成。
ユーザーが対話的に内容を埋めるか、create-design-doc スキルで生成する。

### dependabot.yml

```yaml
version: 2
updates:
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
```

技術スタック固有のエコシステム（gomod, npm等）は init-backend / init-frontend が追記する。

## 前提条件

- `git` がインストール済み
- `gh` CLI がインストール・認証済み
- `mise` がインストール済み

## 次のステップ

init-project 完了後、技術スタックに応じて以下を実行：

```bash
/init-go-backend      # Go バックエンドを追加
/init-react-frontend  # React フロントエンドを追加
```

## 関連スキル

- `init-go-backend`: Go バックエンドセットアップ
- `init-react-frontend`: React フロントエンドセットアップ
- `create-design-doc`: 設計ドキュメント作成
- `create-feature-brief`: 要件定義作成
