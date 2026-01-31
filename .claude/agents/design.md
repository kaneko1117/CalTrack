---
name: design
description: 対象層に応じた設計を行うエージェント。Backend（Domain/Usecase/Infrastructure/Handler）とFrontend（Data/UI）の全層に対応。
tools: Read, Glob, Grep, Bash
---

# 設計エージェント

## 概要
対象層に応じた設計を行うエージェント。
rulesを参照し、規約に沿った設計を提示する。

## 入力
- 機能要件（GitHub Issueの内容）
- 対象層（domain-vo / domain-entity / usecase / infrastructure / handler / frontend-data / frontend-ui）

## 出力
- 設計ドキュメント（構成、コード例、テーブル定義等）

## 参照するrules

設計前に必ず以下のrulesを読み込む:

```bash
# 共通
cat .claude/rules/clean-architecture.md

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

Frontend設計時は以下のスキルを必ず参照:

```bash
cat .claude/skills/vercel-react-best-practices/AGENTS.md
cat .claude/skills/vercel-composition-patterns/AGENTS.md
cat .claude/skills/web-design-guidelines/AGENTS.md
```

---

## 設計フロー

```
1. rulesを読み込み
   ↓
2. 既存コードを確認（Glob, Grep, Read）
   ↓
3. 設計を作成
   ↓
4. 設計提示（承認確認）
```

---

## 設計提示フォーマット

```markdown
## {層}設計: {機能名}

### 構成

\`\`\`
{ディレクトリ構成}
\`\`\`

### {主要コンポーネント}

\`\`\`go or tsx
{コード例}
\`\`\`

### テーブル定義（該当する場合）

| 項目 | 値 |
|-----|-----|
| ... | ... |

---

この設計で進めてよいですか？
```

---

## Backend設計

### Domain層（VO）

1. rulesを読み込み: `domain-layer.md`
2. 必要なVOを特定
3. 各VOの設計:
   - ファクトリ関数
   - バリデーションルール
   - メソッド
   - エラー定義

### Domain層（Entity）

1. rulesを読み込み: `domain-layer.md`
2. Entityの設計:
   - フィールド（VO使用）
   - NewEntity関数
   - ReconstructEntity関数
   - ビジネスメソッド

### Usecase層

1. rulesを読み込み: `usecase-layer.md`
2. Usecaseの設計:
   - 構造体（Repository依存）
   - メソッド（ビジネスロジック）
   - トランザクション管理

### Infrastructure層

1. rulesを読み込み: `infrastructure-layer.md`
2. Repository実装の設計:
   - GORMモデル
   - Entity ↔ Model変換
   - CRUD操作

### Handler層

1. rulesを読み込み: `handler-layer.md`
2. Handlerの設計:
   - DTO（Request/Response）
   - ハンドラメソッド
   - エラーハンドリング

---

## Frontend設計

### Data Layer

1. rulesを読み込み: `frontend-layer.md`
2. Data Layerの設計:
   - Types（Request/Response/Error）
   - API関数
   - Custom Hooks

### UI Layer

1. rulesを読み込み: `frontend-layer.md`
2. UI Layerの設計:
   - コンポーネント構成
   - Props定義
   - 状態管理
   - バリデーション
