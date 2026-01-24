# <機能名> 設計書

## 1. 概要

`<機能の目的と背景を1〜3文で記述>`

### 対象エンティティ / リソース

| 項目 | 内容 |
|------|------|
| `<対象1>` | `<処理内容>` |
| `<対象2>` | `<処理内容>` |

### 前提条件

- `<ドメインルールや事前状態の制約>`
- `<外部サービスの要件>`

---

## 2. レイヤ構成

`<プロジェクトのアーキテクチャに合わせたレイヤ図を記述>`

---

## 3. Domain 層

### 3.1 インターフェース

```python
class <ServiceInterface>(ABC):
    """<サービスの説明>"""

    @abstractmethod
    def <method_name>(self, <params>) -> <ResultType>:
        """<メソッドの説明>"""
        pass
```

### 3.2 型定義

```python
class <ResultType>(TypedDict):
    """<結果の説明>"""
    <field_1>: str
    <field_2>: str
```

### 3.3 例外

| 例外 | 基底 | エラーコード | HTTP |
|------|------|-------------|------|
| `<ServiceError>` | `<BaseException>` | `<ERROR_CODE>` | `<status>` |

### 3.4 Config

```python
@dataclass(frozen=True)
class <ServiceConfig>:
    """<サービス名> 設定。"""
    <field_1>: str
    <field_2>: str
```

---

## 4. Infrastructure 層

### 4.1 認証

```python
def get_<service>_credentials(config: <ServiceConfig>) -> <CredentialType>:
    """<認証方式の説明>"""
```

### 4.2 クライアント生成

```python
def get_<service>_client(config: <ServiceConfig>) -> <ClientType>:
    """<クライアントの説明>"""
```

### 4.3 インターフェース実装

```python
class <ServiceImpl>(<ServiceInterface>):
    def __init__(self, config: <ServiceConfig>):
        self._config = config
        self._client = get_<service>_client(config)

    def <method_name>(self, ...) -> <ResultType>:
        """
        <実装の概要>
        """
```

**リトライ戦略**: `<リトライ対象のエラーと方式>`

---

## 5. Application 層

### 5.1 UseCase: `<UseCaseName>`

```python
class <UseCaseName>:
    def __init__(self, <dependencies>): ...

    def execute(self, <params>) -> <UseCaseResult>:
        """
        1. <ステップ1>
        2. <ステップ2>
        3. <ステップ3>
        """
```

### 5.2 Orchestrator（複数UseCase統合時）

```python
class <OrchestratorName>:
    def execute(self, <params>) -> <OrchestratorResult>:
        """
        複数ステップを順次実行:
        1. <ステップ1>
        2. <ステップ2>
        """
```

---

## 6. API エンドポイント

| メソッド | パス | 説明 |
|----------|------|------|
| `POST` | `/<resource>` | `<操作の説明>` |
| `GET` | `/<resource>/{id}` | `<操作の説明>` |

---

## 7. テスト戦略

### 7.1 テストファイル構成

```
tests/
├── domain/
│   └── test_<domain_logic>.py
├── infrastructure/
│   └── test_<service>.py
└── application/
    └── test_<use_case>.py
```

### 7.2 テスト観点

| レイヤ | 観点 |
|--------|------|
| Domain | ビジネスルール検証 |
| Infrastructure | 外部API連携、エラー変換 |
| Application | 正常フロー、前提条件違反 |

---

## 8. 実装順序と依存関係

```
Phase 1: Domain + Config 基盤
       │
       ▼
Phase 2: Infrastructure + UseCase（並列可）
       │
       ▼
Phase 3: API エンドポイント
       │
       ▼
Phase 4: 統合テスト
```

---

## 9. セキュリティ考慮事項

- `<認証・認可に関する考慮>`
- `<機密情報の取り扱い>`

---

## 10. 環境変数

```bash
# <サービス名>
<ENV_VAR_1>=
<ENV_VAR_2>=
```

---

## 11. 依存パッケージ

```
<package>~=<version>
```
