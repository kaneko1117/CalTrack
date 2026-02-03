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
