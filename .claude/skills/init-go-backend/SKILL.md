---
name: init-go-backend
description: Set up Go backend with clean architecture layers, linter, and formatter. Run after init-project.
allowed-tools: Bash, Write, Read, Edit
---

# Init Go Backend

Go バックエンドのセットアップを行います。

## 概要

このスキルは以下を実行します：

```
1. 情報収集（対話）
   - Go モジュールパス（例: github.com/user/project）
   - 追加パッケージ（デフォルト: chi, yaml.v3, fsnotify, websocket）
   - レイヤー構成（デフォルト: clean architecture）
        ↓
2. Go モジュール初期化
   - go mod init
   - 依存パッケージ追加
        ↓
3. レイヤーディレクトリ作成
   - internal/domain/
   - internal/usecase/
   - internal/handler/
   - internal/infra/
   - cmd/<app-name>/
        ↓
4. リンター設定
   - .golangci.yml（depguard でレイヤー依存ルール強制）
        ↓
5. mise.toml 更新
   - Go ツール追加（golangci-lint, goimports）
   - タスク追加（build, dev, test, test:coverage, fmt, lint）
        ↓
6. Pre-commit hook 更新
   - lint-staged に Go 設定追加（goimports）
   - pre-commit に golangci-lint 追加
        ↓
7. dependabot 更新
   - gomod エコシステム追加
        ↓
8. CLAUDE.md 更新
   - レイヤールール追記
   - Go コマンド追記
```

## 使用方法

```bash
/init-go-backend
```

## 作成・更新されるファイル

### 新規作成

```
├── cmd/<app-name>/
│   └── main.go              (placeholder)
├── internal/
│   ├── domain/
│   │   ├── board.go         (空パッケージ宣言)
│   │   ├── card.go
│   │   └── repository.go
│   ├── usecase/
│   ├── handler/
│   └── infra/
│       ├── yaml/
│       └── watcher/
├── go.mod
├── go.sum
└── .golangci.yml
```

### 更新

```
├── mise.toml               (Go ツール + タスク追加)
├── .husky/pre-commit       (golangci-lint 追加)
├── package.json            (lint-staged Go 設定追加)
├── .github/dependabot.yml  (gomod 追加)
└── CLAUDE.md               (レイヤールール追記)
```

## .golangci.yml の depguard ルール

レイヤー間の依存方向を自動強制：

```
handler → usecase → domain ← infra
```

- domain: 他の internal パッケージを import 禁止
- usecase: handler, infra を import 禁止
- handler: infra を import 禁止

## 前提条件

- `init-project` が実行済み
- `mise` がインストール済み（Go は mise 経由でインストール）

## カスタマイズ

### レイヤー構成の変更

デフォルトは Clean Architecture（domain, usecase, handler, infra）。
別の構成が必要な場合は対話時に指定可能：

- **Flat**: `internal/` 直下にすべて配置
- **Hexagonal**: ports/adapters 構成
- **Custom**: ユーザー指定

### 依存パッケージ

デフォルトパッケージ：
- `github.com/go-chi/chi/v5` — HTTP ルーター
- `gopkg.in/yaml.v3` — YAML
- `github.com/fsnotify/fsnotify` — ファイル監視
- `github.com/gorilla/websocket` — WebSocket

対話時に追加・削除可能。

## 関連スキル

- `init-project`: プロジェクト基盤（先に実行）
- `init-react-frontend`: フロントエンドセットアップ
