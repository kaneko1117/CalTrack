---
paths:
  - "frontend/src/**/*.ts"
  - "frontend/src/**/*.tsx"
---

# Frontend層規則

## 必須参照スキル

Frontend作業時は必ず以下のスキルを参照すること:

```bash
cat .claude/skills/vercel-react-best-practices/AGENTS.md
cat .claude/skills/vercel-composition-patterns/AGENTS.md
cat .claude/skills/web-design-guidelines/AGENTS.md
```

| スキル | 用途 |
|-------|------|
| vercel-react-best-practices | パフォーマンス最適化（40+ルール） |
| vercel-composition-patterns | コンポーネント構成パターン |
| web-design-guidelines | Webデザインガイドライン |

---

## ディレクトリ構成

```
frontend/src/
├── domain/             # DDD Domain層
│   ├── shared/         # 共通型（Result, DomainError）
│   ├── valueObjects/   # Value Object
│   └── entities/       # Entity
├── features/           # 機能単位
│   └── {feature}/
│       ├── api/
│       │   └── index.ts
│       ├── hooks/
│       │   ├── use{Feature}.ts  # 型もこの中で定義
│       │   └── index.ts
│       ├── components/
│       │   ├── {Component}.tsx  # 型もこの中で定義
│       │   ├── {Component}.test.tsx
│       │   └── index.ts
│       └── index.ts
├── components/ui/      # shadcn/ui コンポーネント
├── pages/              # ページコンポーネント
├── routes/             # ルーティング
├── hooks/              # 共通フック
├── lib/                # ユーティリティ
│   └── api.ts          # Axiosインスタンス
├── types/              # 共通型定義（複数feature横断のみ）
└── test/               # テスト設定
```

**重要: `types/index.ts` は作成しない。型はコンポーネントやhooks内で定義する。**

---

## Domain層（domain/）

フロントエンドでもDDDを採用し、バリデーションロジックをドメイン層に集約する。

### 構成

```
domain/
├── shared/           # Result型、DomainError型
├── valueObjects/     # Value Object
└── entities/         # Entity
```

### 設計思想

| 項目 | 内容 |
|------|------|
| Result型 | 成功/失敗を型安全に表現（Ok/Err） |
| DomainError型 | エラーコード + メッセージを保持 |
| Value Object | 不変、ファクトリ関数で生成、バリデーション内包 |
| Entity | 複数VOを組み合わせ、エラーを集約して返す |

### VOの特徴

- 不変（`Object.freeze`）
- `new{VO}()` でバリデーション付き生成
- `Result<VO, DomainError>` を返す
- エラーはcode + messageを持つ（hooksから単独呼び出し可能）
- `equals()` で等価性比較

### 禁止事項

- 外部ライブラリの使用禁止（純粋なTypeScriptのみ）

---

## Frontend Layer（api / hooks / components）

### 型定義ルール

**`types/index.ts` は作成しない。型はコンポーネントやhooks内で定義する。**

```typescript
// ❌ 禁止: types/index.ts を作成
// features/statistics/types/index.ts

// ✅ 正しい: hooks内で型定義
// features/statistics/hooks/useStatistics.ts
export type StatisticsResponse = {
  period: string;
  targetCalories: number;
  // ...
};

export function useStatistics(period: Period) {
  // ...
}

// ✅ 正しい: コンポーネント内で型定義
// features/statistics/components/CalorieChart.tsx
export type CalorieChartProps = {
  data: DailyStatistic[];
  targetCalories: number;
};

export function CalorieChart({ data, targetCalories }: CalorieChartProps) {
  // ...
}
```

**`interface` ではなく `type` を使用する**

```typescript
/** エラーコード定数 */
export const ERROR_CODE_INTERNAL_ERROR: ErrorCode = "INTERNAL_ERROR";

/** エラーメッセージ定数 */
export const ERROR_MESSAGE_UNEXPECTED = "予期しないエラーが発生しました";
```

### API

```typescript
/**
 * {機能説明}
 * @param request - リクエストデータ
 * @returns Promise<{Feature}Response>
 * @throws ApiError
 */
export async function {featureName}(
  request: {Feature}Request
): Promise<{Feature}Response> {
  // 実装
}
```

### Hooks

```typescript
/**
 * {フック説明}
 * @returns { action, isLoading, error, isSuccess, reset }
 */
export function use{Feature}(): Use{Feature}Return {
  // useState, useCallback で状態管理
}
```

### コンポーネント設計

| 種別 | 特徴 |
|-----|------|
| Container | Hooks使用、ロジック担当 |
| Presentational | propsのみ、再利用可能 |

### 命名規則

- コンポーネント: PascalCase（例: `RegisterForm`）
- ファイル: PascalCase（例: `RegisterForm.tsx`）
- テスト: `{Component}.test.tsx`

### フォームコンポーネント

```typescript
type {Form}Props = {
  onSuccess?: () => void;
};

type FormState = {
  // フォームフィールド
};

type FormErrors = {
  // バリデーションエラー
};

function validate(form: FormState): FormErrors {
  // バリデーション実装
}
```

---

## 共通ルール

### 型定義

- **`interface` ではなく `type` を使用**
- `any` 禁止
- 明示的な型定義

```typescript
// NG
interface UserProps {
  name: string;
}

// OK
type UserProps = {
  name: string;
};
```

### コメント

- **日本語で記述**
- JSDocスタイル使用

```typescript
/**
 * ユーザー登録フック
 * ローディング状態、エラー状態、成功状態を管理
 */
```

### 定数化

- エラーコード、エラーメッセージは定数化
- マジックナンバー禁止

---

## テスト

- Vitest + React Testing Library
- ファイル: `{Component}.test.tsx` または `{hook}.test.ts`
- カバレッジ: 正常系・異常系・境界値

---

## Storybook

**コンポーネント実装時は必ずStoryを作成すること**

- ファイル: `{Component}.stories.tsx`
- 配置: コンポーネントと同じディレクトリ
- 必須Story: Default, Loading（該当時）, Empty（該当時）

---

## shadcn/ui

- `components/ui/` に配置
- 公式スタイル準拠
- 必要に応じて追加（`npx shadcn-ui@latest add {component}`）

---

## Feature間依存ルール

**feature間の直接依存は禁止**

```
features/
├── common/        # 共通機能（他featureから参照可）
├── records/       # records は common のみ参照可
├── statistics/    # statistics は common のみ参照可
└── ...
```

| 参照元 | 参照可能 | 参照禁止 |
|-------|---------|---------|
| records | common | statistics, その他feature |
| statistics | common | records, その他feature |
| その他 | common | 他のfeature |

**例外: common は全featureから参照可能**

---

## データ要素アニメーションルール

**データ要素を表示する際は以下のアニメーションを必須とする**

| 要素 | 必須アニメーション |
|------|------------------|
| カード | フェードイン（`opacity-0 animate-fade-in-up`） |
| 数値 | カウントアップ（`useCountUp` フック使用） |
| ローディング | スケルトン（`Skeleton` コンポーネント使用） |

**フェードインの遅延**
- 複数要素がある場合は `animation-delay-{n}` で順次表示
- 例: `animation-delay-100`, `animation-delay-200`, `animation-delay-300`
