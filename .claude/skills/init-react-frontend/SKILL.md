---
name: init-react-frontend
description: Set up React frontend with Vite, TypeScript, ESLint, Prettier, and dev proxy. Run after init-project.
allowed-tools: Bash, Write, Read, Edit
---

# Init React Frontend

React フロントエンドのセットアップを行います。

## 概要

このスキルは以下を実行します：

```
1. 情報収集（対話）
   - 追加パッケージ（例: @dnd-kit, react-router, zustand）
   - API proxy 先ポート（デフォルト: 8080）
   - ディレクトリ名（デフォルト: web）
        ↓
2. Vite プロジェクト初期化
   - npm create vite -- --template react-ts
   - npm install
        ↓
3. 追加パッケージインストール
   - ユーザー指定パッケージ
   - prettier, eslint-config-prettier（共通）
        ↓
4. ESLint 設定
   - eslint-config-prettier 統合
   - React Hooks ルール有効化
        ↓
5. Prettier 設定
   - .prettierrc 作成
        ↓
6. Vite 設定更新
   - API proxy 設定（/api → バックエンド）
   - WebSocket proxy 設定（/ws → バックエンド）
        ↓
7. mise.toml 更新
   - Node ツール追加
   - フロントエンドタスク追加（dev:front, fmt, lint）
        ↓
8. Pre-commit hook 更新
   - lint-staged に TS/TSX 設定追加（prettier + eslint）
        ↓
9. dependabot 更新
   - npm エコシステム追加
        ↓
10. CLAUDE.md 更新
    - フロントエンドコマンド追記
```

## 使用方法

```bash
/init-react-frontend
```

## 作成・更新されるファイル

### 新規作成

```
web/
├── index.html
├── package.json
├── package-lock.json
├── tsconfig.json
├── tsconfig.app.json
├── tsconfig.node.json
├── vite.config.ts
├── eslint.config.js
├── .prettierrc
├── public/
└── src/
    ├── App.tsx
    ├── main.tsx
    ├── App.css
    └── index.css
```

### 更新

```
├── mise.toml               (Node + フロントタスク追加)
├── package.json            (lint-staged TS 設定追加)
├── .github/dependabot.yml  (npm 追加)
└── CLAUDE.md               (フロントコマンド追記)
```

## 設定内容

### Prettier (.prettierrc)

```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "all"
}
```

### ESLint

- TypeScript strict
- React Hooks ルール
- eslint-config-prettier で整形ルールを無効化

### Vite Dev Proxy

```typescript
server: {
  proxy: {
    '/api': 'http://localhost:<port>',
    '/ws': { target: 'http://localhost:<port>', ws: true }
  }
}
```

### TypeScript

- `strict: true`
- `noUnusedLocals: true`
- `noUnusedParameters: true`

## 前提条件

- `init-project` が実行済み
- `mise` がインストール済み（Node は mise 経由でインストール）

## カスタマイズ

### よく使う追加パッケージ

| パッケージ | 用途 |
|-----------|------|
| `@dnd-kit/core @dnd-kit/sortable` | ドラッグ&ドロップ |
| `react-router-dom` | ルーティング |
| `zustand` | 状態管理 |
| `@tanstack/react-query` | サーバー状態管理 |
| `clsx` | クラス名結合 |

対話時に指定すると自動インストール。

### CSS方式

デフォルトは CSS Modules（Vite 標準対応）。
以下も対話時に選択可能：

- **Tailwind CSS**: `tailwindcss` + PostCSS 設定追加
- **CSS-in-JS**: `styled-components` 等
- **Plain CSS**: そのまま

## 関連スキル

- `init-project`: プロジェクト基盤（先に実行）
- `init-go-backend`: バックエンドセットアップ
