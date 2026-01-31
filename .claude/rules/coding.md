---
paths:
  - "backend/**/*.go"
  - "frontend/src/**/*.ts"
  - "frontend/src/**/*.tsx"
---

# コーディング規約

## 共通ルール

### コメント

- **日本語で記述**
- 簡潔で明確に
- 「何をしているか」ではなく「なぜそうしているか」を書く

### 定数化

- マジックナンバー禁止
- エラーコード、エラーメッセージは定数化
- 設定値は環境変数または定数ファイルに

### ヘルパー関数

- 共通処理はヘルパー関数に抽出
- 再利用可能な形で設計

---

## Backend (Go)

### フォーマット

```bash
go fmt ./...
```

### 命名規則

| 対象 | 規則 | 例 |
|-----|------|-----|
| パッケージ | 小文字、短く | `user`, `auth` |
| 構造体 | PascalCase | `UserHandler` |
| メソッド | PascalCase | `Register` |
| 変数 | camelCase | `userRepo` |
| 定数 | PascalCase または UPPER_SNAKE | `ErrNotFound` |

### ファイル構成

- ドメイン単位で1ファイル
- 例: `usecase/user.go`, `handler/user/handler.go`

### エラーハンドリング

```go
if err != nil {
    return err
}
```

- エラーは早期リターン
- ラップする場合は `fmt.Errorf("context: %w", err)`

---

## Frontend (TypeScript/React)

### フォーマット

```bash
npm run lint
npm run format
```

### 命名規則

| 対象 | 規則 | 例 |
|-----|------|-----|
| コンポーネント | PascalCase | `RegisterForm` |
| 関数 | camelCase | `registerUser` |
| 変数 | camelCase | `isLoading` |
| 定数 | UPPER_SNAKE | `ERROR_CODE_INTERNAL_ERROR` |
| 型 | PascalCase | `RegisterUserRequest` |
| フック | use + PascalCase | `useRegisterUser` |

### 型定義

- **`interface` ではなく `type` を使用**
- `any` 禁止
- 明示的な型定義を使用
- `unknown` を使い、型ガードで絞り込む

```typescript
// NG
interface UserProps {
  name: string;
}
const data: any = response.data;

// OK
type UserProps = {
  name: string;
};
const data: UserResponse = response.data;
```

### インポート順序

1. 外部ライブラリ
2. 内部モジュール（絶対パス）
3. 相対パス

```typescript
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { useAuth } from "../hooks";
```

---

## Git

### コミットメッセージ

```
{type}({scope}): {subject}

{body}

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
```

| type | 用途 |
|-----|------|
| feat | 新機能 |
| fix | バグ修正 |
| docs | ドキュメント |
| refactor | リファクタリング |
| test | テスト |
| chore | その他 |

### ブランチ命名

- `feat/{feature-name}`
- `fix/{bug-name}`
- `refactor/{target}`
