# アーキテクチャ規則

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

## 集計ロジックの配置方針

### 原則
- **基本的な集計はSQLで行う**（Repository Interface経由）
- **複雑化する場合はアプリケーション側（Domainのメソッド）で実施する**

### 判断基準

| 集計の種類 | 配置先 | 例 |
|-----------|--------|-----|
| 単純な合計・平均・グルーピング | SQL（Repository） | 日別カロリー合計、期間内の平均値 |
| 条件付き集計（WHERE + GROUP BY程度） | SQL（Repository） | ユーザー別・期間別の集計 |
| 複数テーブルの集計結果を組み合わせる計算 | Domain（Entityメソッド） | PFCバランスの評価、目標との差分計算 |
| ビジネスルールに基づく分類・判定を伴う集計 | Domain（Entityメソッド） | 栄養素の過不足判定、スコアリング |

### 各層での扱い

- **Repository Interface**（`domain/repository/`）: 集計メソッドと結果型を定義する。結果型にはVOを使用する
- **Infrastructure**（Repository実装）: SQLで集計を実行する
- **Usecase**: Repositoryから集計結果を取得し、必要に応じてDomainのメソッドで追加の計算を行う
- **Domain**: 複雑な計算・判定ロジックをEntityまたはVOのメソッドとして持つ
