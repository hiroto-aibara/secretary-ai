---
name: init-serena
description: Set up Serena MCP for semantic code operations. Adds MCP server config and gitignore entry.
allowed-tools: Bash, Read, Edit
---

# Init Serena

Serena MCP をプロジェクトにセットアップします。

## 概要

Serena はセマンティックなコード操作を提供する MCP サーバーです。
シンボル検索・リネーム・参照検索など、LSP ベースの操作が可能になります。

```
1. 前提条件の確認
   - uv がインストールされているか
        ↓
2. Serena MCP をプロジェクトに追加
   - claude mcp add コマンド実行
        ↓
3. .gitignore に .serena/ を追加
        ↓
4. 完了メッセージ
```

## 使用方法

```bash
/init-serena
```

## 実行手順

### 1. 前提条件の確認

```bash
which uv
```

`uv` が見つからない場合は、インストール方法を案内して終了：

```bash
# macOS / Linux
curl -LsSf https://astral.sh/uv/install.sh | sh
```

### 2. Serena MCP を追加

```bash
claude mcp add serena -- uvx --from git+https://github.com/oraios/serena serena start-mcp-server --context claude-code --project "$(pwd)"
```

**パラメータ説明：**
- `--context claude-code`: Claude Code と重複するツールを無効化
- `--project "$(pwd)"`: 現在のプロジェクトをアクティベート

### 3. .gitignore に追加

`.gitignore` に以下を追加（まだ存在しない場合）：

```
# Serena
.serena/
```

### 4. 完了メッセージ

以下を表示：

```
Serena MCP をセットアップしました。

次回の Claude Code セッションから以下のツールが利用可能です：
- get_symbols_overview: ファイル内のシンボル一覧
- find_symbol: シンボル検索
- find_referencing_symbols: 参照検索
- rename_symbol: シンボルリネーム（全参照を自動更新）
- replace_symbol_body: シンボル単位の編集

Note: メモリ機能は使用しません。
コンテキストは CLAUDE.md + docs/ を Source of Truth とします。
```

## 運用方針

### メモリについて

Serena にはメモリ機能がありますが、このプロジェクトでは**使用しません**。

理由：
- `CLAUDE.md` + `docs/` が Source of Truth
- 二重管理によるコンテキスト汚染を防ぐ

### Serena の主な用途

| 機能 | 用途 |
|------|------|
| `get_symbols_overview` | ファイル構造の把握 |
| `find_symbol` | クラス・関数の検索 |
| `find_referencing_symbols` | 「この関数を呼んでいる箇所」の検索 |
| `rename_symbol` | リファクタリング（全参照を自動更新） |
| `replace_symbol_body` | シンボル単位の編集 |

## 前提条件

- `uv` パッケージマネージャーがインストール済み
- `claude` CLI がインストール済み

## 関連リンク

- [Serena GitHub](https://github.com/oraios/serena)
- [Serena User Guide](https://oraios.github.io/serena/)
