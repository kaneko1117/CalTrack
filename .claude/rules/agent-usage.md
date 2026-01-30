---
alwaysApply: true
---

# カスタムエージェント優先使用ルール

## 概要
このプロジェクトでは `.claude/agents/` にカスタムエージェントが定義されている。
機能開発時は、これらのエージェントを優先的に使用すること。

## カスタムエージェント一覧

### 司令塔
| エージェント | 用途 |
|-------------|------|
| `orchestrator` | 仕様書から設計→実装→PR作成の全体フローを管理 |

### 設計エージェント
| エージェント | 用途 |
|-------------|------|
| `domain-design` | Backend Domain層（VO + Entity）の詳細設計 |
| `usecase-design` | Backend Usecase層の詳細設計 |
| `infrastructure-design` | Backend Infrastructure層の詳細設計 |
| `handler-design` | Backend Handler層の詳細設計 |
| `frontend-data-design` | Frontend Data Layer（types + api + hooks）の詳細設計 |
| `frontend-ui-design` | Frontend UI Layer（components）の詳細設計 |

### 実装エージェント
| エージェント | 用途 |
|-------------|------|
| `impl` | 設計を受け取りコード実装（Backend/Frontend両対応） |
| `test-pr` | Build/Test実行、PR作成（Backend/Frontend両対応） |

## 使用ルール

### 機能開発時
1. **新機能の実装依頼** → `orchestrator` を使用
2. **特定層のみの設計依頼** → 該当する設計エージェントを使用
3. **実装のみの依頼** → `impl` を使用
4. **テスト/PR作成のみ** → `test-pr` を使用

### エージェント使用の判断基準

| 依頼内容 | 使用エージェント |
|---------|----------------|
| 「〇〇機能を実装して」 | `orchestrator` |
| 「Domain層を設計して」 | `domain-design` |
| 「Usecaseを設計して」 | `usecase-design` |
| 「この設計を実装して」 | `impl` |
| 「テストしてPR作って」 | `test-pr` |
| 「フロントのhooksを設計して」 | `frontend-data-design` |
| 「コンポーネントを設計して」 | `frontend-ui-design` |

### 禁止事項
- 設計エージェントを使わず直接コードを書くこと（簡単な修正を除く）
- orchestrator のフローを無視して層を飛ばすこと
- PR作成時に設計書をスキップすること

## エージェントの呼び出し方

```
Task tool で subagent_type=general-purpose を使用し、
エージェント定義ファイル（.claude/agents/{name}.md）を読み込んで
その指示に従って実行する。
```
