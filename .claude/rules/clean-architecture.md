# クリーンアーキテクチャ規則

## 依存関係の方向

```
Handler → Usecase → Domain
              ↓
        Infrastructure
```

- 依存は内側（Domain）に向かう
- Domain層は他の層に依存しない
- Infrastructure層はUsecase層のInterfaceを実装する

## 層の責務

| 層 | 責務 | 禁止事項 |
|----|------|---------|
| Domain | ビジネスルール、Entity、VO | フレームワーク依存、DB依存 |
| Usecase | アプリケーションロジック | HTTP依存、DB直接アクセス |
| Infrastructure | 技術的実装（DB、外部API） | ビジネスロジック |
| Handler | HTTP処理、入出力変換 | ビジネスロジック、DB直接アクセス |

## import規則

- `domain/` は他の層をimportしない
- `usecase/` は `domain/` のみimport可
- `infrastructure/` は `domain/`, `usecase/` をimport可
- `handler/` は `domain/`, `usecase/` をimport可（`infrastructure/` は不可）

## 集計・計算ロジックの配置

**集計・計算処理はアプリケーション層（Domain/Usecase）で行う。SQLでは行わない。**

### 理由

1. **ビジネスロジックの集約**: 集計ロジック（合計、平均、カウント等）はビジネスルールの一部
2. **テスタビリティ**: 単体テストが書きやすい
3. **一貫性**: ロジックがDomain/Usecase層に統一される

### 実装パターン

```go
// ✅ 正しい: Repositoryはデータ取得のみ
records, _ := repo.FindByUserIDAndDateRange(ctx, userID, start, end)

// ✅ 正しい: 集計はUsecase/Entity層で行う
totalCalories := 0
for _, record := range records {
    totalCalories += record.TotalCalories()
}
```

```go
// ❌ 禁止: RepositoryでSQLの集計関数を使用
func (r *repo) GetTotalCalories(ctx, userID) (int, error) {
    // SELECT SUM(calories) FROM records WHERE user_id = ?
}
```

### 例外

- パフォーマンス問題が明確に発生した場合のみ、SQLでの集計を検討
- その場合も、Usecase層で集計ロジックを定義し、Repositoryは最適化された実装として扱う
