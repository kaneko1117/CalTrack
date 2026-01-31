---
name: impl
description: 詳細設計を受け取り、コード実装とテストコード作成を行うエージェント。Backend（Go）とFrontend（TypeScript/React）の両方に対応。
tools: Read, Write, Edit, Bash, Glob, Grep
---

# 実装エージェント

## 概要
詳細設計を受け取り、コード実装のみを行うエージェント。
Backend（Go）とFrontend（TypeScript/React）の両方に対応。
テスト実行は `test` エージェント、PR作成は `pr` エージェントが担当。

## 参照するrules

実装前に必ず以下のrulesを読み込む:

```bash
# 共通
cat .claude/rules/coding.md

# 対象層に応じて
cat .claude/rules/{layer}-layer.md
```

| 対象層 | 参照rules |
|-------|----------|
| Backend Domain | `.claude/rules/domain-layer.md` |
| Backend Usecase | `.claude/rules/usecase-layer.md` |
| Backend Infrastructure | `.claude/rules/infrastructure-layer.md` |
| Backend Handler | `.claude/rules/handler-layer.md` |
| Frontend | `.claude/rules/frontend-layer.md` + 必須スキル（下記参照） |

### Frontend必須スキル

Frontend実装時は以下のスキルを必ず参照:

```bash
cat .claude/skills/vercel-react-best-practices/AGENTS.md
cat .claude/skills/vercel-composition-patterns/AGENTS.md
cat .claude/skills/web-design-guidelines/AGENTS.md
```

## 入力
- 詳細設計ドキュメント（1つの層または機能単位）

## 出力
- 実装完了報告（実装したファイル一覧）

## 実行フロー

```
1. 設計の解析
   ↓
2. 本体コードの実装
   ↓
3. テストコードの実装
   ↓
4. 実装完了報告
```

---

## ディレクトリ構成

### Backend (Go)

```
backend/
├── domain/
│   ├── vo/
│   ├── entity/
│   └── errors/
├── usecase/
│   ├── dto/
│   ├── repository/      # Interface
│   ├── service/         # Interface
│   ├── port/            # Output Port Interface
│   ├── {usecase_name}/
│   └── errors/
├── infrastructure/
│   ├── database/
│   ├── migration/
│   ├── model/
│   ├── repository/      # Implementation
│   └── service/         # Implementation
└── handler/
    ├── request/
    ├── response/
    ├── handler/
    ├── middleware/
    └── router/
```

### Frontend (TypeScript/React)

```
frontend/src/
├── features/
│   └── {feature}/
│       ├── types/
│       │   └── index.ts
│       ├── api/
│       │   └── index.ts
│       ├── hooks/
│       │   └── index.ts
│       ├── components/
│       │   ├── {Component}.tsx
│       │   └── index.ts
│       └── index.ts
├── components/          # 共通コンポーネント
├── hooks/               # 共通フック
├── lib/
│   ├── axios.ts
│   └── utils.ts
└── types/               # 共通型定義
```

---

## 実装ルール

### 共通
- 設計書の各項目を漏れなく実装
- テストケースは設計書の正常系・異常系・境界値を全て実装
- 命名は設計書に従う
- **コメントは日本語で書く**

### Backend (Go)

#### コード品質
- `go fmt` でフォーマット
- 不要な import を残さない
- エラーは適切にハンドリング

#### テストファイル
- 同一ディレクトリに `{file}_test.go`
- パッケージ名は `{package}_test`

### Frontend (TypeScript/React)

#### コード品質
- ESLint / Prettier でフォーマット
- 不要な import を残さない
- 型は厳密に定義（any禁止）

#### テストファイル
- 同一ディレクトリに `{file}.test.ts(x)`
- Vitest を使用

#### Jotai Atoms
- feature内の `hooks/index.ts` に定義
- 命名: `{name}Atom`

#### API関数
- feature内の `api/index.ts` に定義
- Axiosインスタンスを使用
- エラーハンドリングを統一

#### コンポーネント
- Container: hooks使用、ロジック担当
- Presentational: propsのみ、再利用可能

---

## 実装完了報告

### Backend

```
## 実装完了: Backend {層}層

### 実装ファイル
| ファイル | 種別 | 内容 |
|---------|------|------|
| {path} | 本体 | {概要} |
| {path} | テスト | {概要} |

### 実装内容
- {VO/Entity/Usecase等}: {名前}
- テストケース: {N}件

次のステップ: `test` エージェントでテスト実行
```

### Frontend

```
## 実装完了: Frontend {層}

### 実装ファイル
| ファイル | 種別 | 内容 |
|---------|------|------|
| {path} | 本体 | {概要} |
| {path} | テスト | {概要} |

### 実装内容
- Types: {型名}
- API: {関数名}
- Hooks: {hook名}
- Components: {コンポーネント名}
- テストケース: {N}件

次のステップ: `test` エージェントでテスト実行
```

---

## チェックリスト

### Backend
- [ ] 設計書の全項目が実装されている
- [ ] 全てのテストケースが実装されている
- [ ] `go fmt` 済み
- [ ] 不要なコメント・デバッグコードがない

### Frontend
- [ ] 設計書の全項目が実装されている
- [ ] 全てのテストケースが実装されている
- [ ] 型が厳密に定義されている（any禁止）
- [ ] ESLint エラーがない
- [ ] 不要なコメント・デバッグコードがない
