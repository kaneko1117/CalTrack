---
paths:
  - "backend/**/*.go"
  - "web/src/**/*.ts"
  - "web/src/**/*.tsx"
---

# 共通規約

コーディング規約とテスト規則を統合したルール。

---

## コーディング規約

### 共通ルール

#### コメント
- **日本語で記述**
- 簡潔で明確に
- 「何をしているか」ではなく「なぜそうしているか」を書く

#### 定数化
- マジックナンバー禁止
- エラーコード、エラーメッセージは定数化
- 設定値は環境変数または定数ファイルに

#### ヘルパー関数
- 共通処理はヘルパー関数に抽出
- 再利用可能な形で設計

---

### Backend (Go)

#### フォーマット
```bash
go fmt ./...
```

#### 命名規則
| 対象 | 規則 | 例 |
|-----|------|-----|
| パッケージ | 小文字、短く | `user`, `auth` |
| 構造体 | PascalCase | `UserHandler` |
| メソッド | PascalCase | `Register` |
| 変数 | camelCase | `userRepo` |
| 定数 | PascalCase または UPPER_SNAKE | `ErrNotFound` |

#### ファイル構成
- ドメイン単位で1ファイル
- 例: `usecase/user.go`, `handler/user/handler.go`

#### エラーハンドリング
```go
if err != nil {
    return err
}
```
- エラーは早期リターン
- ラップする場合は `fmt.Errorf("context: %w", err)`

---

### Frontend (TypeScript/React)

#### フォーマット
```bash
npm run lint
npm run format
```

#### 命名規則
| 対象 | 規則 | 例 |
|-----|------|-----|
| コンポーネント | PascalCase | `RegisterForm` |
| 関数 | camelCase | `registerUser` |
| 変数 | camelCase | `isLoading` |
| 定数 | UPPER_SNAKE | `ERROR_CODE_INTERNAL_ERROR` |
| 型 | PascalCase | `RegisterUserRequest` |
| フック | use + PascalCase | `useRegisterUser` |

#### 型定義
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

#### インポート順序
1. 外部ライブラリ
2. 内部モジュール（絶対パス）
3. 相対パス

```typescript
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { useAuth } from "../hooks";
```

---

### Git

#### コミットメッセージ
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

#### ブランチ命名
- `feat/{feature-name}`
- `fix/{bug-name}`
- `refactor/{target}`

---

## テスト規則

### Backend (Go)

#### テストファイル配置
- 同一パッケージ内に `{file}_test.go`
- パッケージ名は `{package}_test`（外部テスト）

#### テスト命名
```
Test{対象}_{メソッド名}
```
例: `TestNewEmail`, `TestRecordHandler_Create`

#### サブテスト（t.Run）
- t.Runを使用してサブテストを定義する
- **テストケース名は日本語**で記述

```go
func TestRecordHandler_Create(t *testing.T) {
    t.Run("正常系_記録が作成される", func(t *testing.T) {
        // ...
    })
    t.Run("異常系_認証エラー", func(t *testing.T) {
        // ...
    })
}
```

#### テーブル駆動テスト

**使うべき場面:**
- 同じロジックを異なる入力でテストする場合
- 境界値テスト
- 入力→出力が単純なマッピングの場合

**使わないべき場面:**
- テストごとに検証ロジックが異なる場合
- 条件分岐がテスト内で多くなる場合
- 正常系の詳細な検証（個別テストの方が読みやすい）

```go
// 良い例: シンプルなバリデーションエラーテスト
tests := []struct {
    name    string
    input   string
    wantErr error
}{
    {"invalid email", "invalid", ErrInvalidEmailFormat},
    {"empty email", "", ErrEmailRequired},
}

// 悪い例: 条件分岐が多いテーブル駆動
// → 個別のテスト関数に分ける
```

---

### Frontend (TypeScript/React)

#### テストファイル配置
- `{Component}.test.tsx` または `{hook}.test.ts`

#### テスト命名
- **テストケース名は日本語**で記述
- `describe`, `it` の第1引数は日本語

```typescript
describe("RegisterForm", () => {
  it("正常系_フォーム送信でユーザーが登録される", () => {
    // ...
  });
  it("異常系_バリデーションエラーが表示される", () => {
    // ...
  });
});
```

---

### 共通テストルール

#### モック戦略
| 層 | モック対象 | 方法 |
|----|----------|------|
| Domain | なし | 実際の値で検証 |
| Usecase | Repository, Service | Interface実装 |
| Infrastructure | DB | sqlmock または testcontainers |
| Handler | Usecase | Interface実装 |

#### カバレッジ基準
- 正常系: 全パターン
- 異常系: 各エラー条件
- 境界値: 最小/最大の境界

#### 禁止事項
- 外部サービスへの実際の接続
- テスト間の状態共有
- time.Now() の直接使用（注入する）
