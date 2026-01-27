# Task: board-list-management

## 概要

フロントエンドUIにボード・リスト管理機能を追加する。
バックエンドAPIとフロントエンドAPI関数は既に存在するため、**UI実装のみ**が対象。

## 要件

### ボード管理
- [ ] ヘッダーに「+ New Board」ボタンを追加
- [ ] ボード作成モーダル（ID/名前/初期リスト入力）を実装
- [ ] ボード選択UI横に削除ボタンを追加
- [ ] ボード削除時に確認ダイアログを表示

### リスト管理
- [ ] ボード末尾に「+ Add List」カラムを追加
- [ ] リスト追加のインライン入力フォームを実装
- [ ] リストヘッダーをクリックで編集モードにする
- [ ] リストのドラッグ&ドロップによる並び替えを実装

## 技術仕様

- 使用技術: React 18, TypeScript, @dnd-kit/core, CSS Modules
- 関連ファイル:
  - `web/src/App.tsx`
  - `web/src/components/Board.tsx`
  - `web/src/components/List.tsx`
  - `web/src/hooks/useApi.ts`（既存API関数を使用）
- API エンドポイント:
  - `POST /api/boards` - ボード作成
  - `PUT /api/boards/{id}` - ボード更新（リスト追加/編集/並び替え）
  - `DELETE /api/boards/{id}` - ボード削除

## 実装ファイル

| ファイル | 変更種別 | 内容 |
|---------|---------|------|
| `web/src/App.tsx` | 修正 | ボード作成/削除UI追加 |
| `web/src/components/Board.tsx` | 修正 | リスト管理機能追加 |
| `web/src/components/List.tsx` | 修正 | 名前編集・ドラッグ対応 |
| `web/src/components/BoardModal.tsx` | 新規 | ボード作成モーダル |
| `web/src/components/BoardModal.module.css` | 新規 | モーダルスタイル |
| `web/src/components/AddList.tsx` | 新規 | リスト追加コンポーネント |
| `web/src/components/AddList.module.css` | 新規 | リスト追加スタイル |
| `web/src/App.module.css` | 修正 | ヘッダーボタン追加 |

## 受け入れ基準

- [ ] 「+ New Board」ボタンでボードを作成できる
- [ ] ボード削除ボタンでボードを削除できる（確認ダイアログあり）
- [ ] 「+ Add List」でリストを追加できる
- [ ] リスト名をクリックして名前を変更できる
- [ ] リストをドラッグ&ドロップで並び替えできる
- [ ] ESLint/Prettier エラーなし
- [ ] 既存のテストがパス

## 品質基準

- 既存のUIパターン（CardModal等）と統一されたデザイン
- レスポンシブ対応（横スクロール可能）
- エラーハンドリング（API失敗時のフィードバック）

## 検証方法

1. `mise run dev` + `mise run dev:front` でアプリ起動
2. http://localhost:5173/ で動作確認:
   - 「+ New Board」でボード作成
   - ボード削除ボタンで削除
   - 「+ Add List」でリスト追加
   - リスト名クリックで名前変更
   - リストをドラッグで並び替え
3. `mise run lint` でリントチェック
4. `mise run test` でテスト実行

## 参考リソース

- 関連計画: `/Users/hiroto_aibara/.claude/plans/reactive-frolicking-dijkstra.md`
- 既存モーダル実装: `web/src/components/CardModal.tsx`
- dnd-kit ドキュメント: https://dndkit.com/
