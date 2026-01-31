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
├── features/           # 機能単位
│   └── {feature}/
│       ├── types/
│       │   └── index.ts
│       ├── api/
│       │   └── index.ts
│       ├── hooks/
│       │   └── index.ts
│       ├── components/
│       │   ├── {Component}.tsx
│       │   ├── {Component}.test.tsx
│       │   └── index.ts
│       └── index.ts
├── components/ui/      # shadcn/ui コンポーネント
├── pages/              # ページコンポーネント
├── routes/             # ルーティング
├── hooks/              # 共通フック
├── lib/                # ユーティリティ
│   └── api.ts          # Axiosインスタンス
├── types/              # 共通型定義
└── test/               # テスト設定
```

---

## Data Layer（types / api / hooks）

### Types

**`interface` ではなく `type` を使用する**

```typescript
/** リクエスト型 */
export type {Feature}Request = {
  // フィールド定義
};

/** レスポンス型 */
export type {Feature}Response = {
  // フィールド定義
};

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

---

## UI Layer（components）

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

## shadcn/ui

- `components/ui/` に配置
- 公式スタイル準拠
- 必要に応じて追加（`npx shadcn-ui@latest add {component}`）
