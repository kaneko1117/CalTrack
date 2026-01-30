---
paths:
  - "backend/**/*_test.go"
---

# テスト規則

## テストファイル配置

- 同一パッケージ内に `{file}_test.go`
- パッケージ名は `{package}_test`（外部テスト）

## テスト命名

```
Test{対象}_{条件}_{期待結果}
```

例:
- `TestNewEmail_ValidFormat_ReturnsEmail`
- `TestUser_ChangePassword_UpdatesPassword`

## テーブル駆動テスト

境界値テストはテーブル駆動で:

```go
tests := []struct {
    name    string
    input   string
    want    Type
    wantErr error
}{
    {"valid case", "input", expected, nil},
    {"invalid case", "bad", nil, ErrInvalid},
}
```

## モック戦略

| 層 | モック対象 | 方法 |
|----|----------|------|
| Domain | なし | 実際の値で検証 |
| Usecase | Repository, Service | Interface実装 |
| Infrastructure | DB | sqlmock または testcontainers |
| Handler | Usecase | Interface実装 |

## カバレッジ基準

- 正常系: 全パターン
- 異常系: 各エラー条件
- 境界値: 最小/最大の境界

## 禁止事項

- 外部サービスへの実際の接続
- テスト間の状態共有
- time.Now() の直接使用（注入する）
