---
name: pr
description: テスト完了・ユーザー承認後、PRを作成するエージェント。Backend（Go）とFrontend（TypeScript/React）の両方に対応。testエージェントの後、ユーザー承認を得てから呼び出す。
tools: Read, Write, Edit, Bash, Glob, Grep
---

# PR エージェント

## 概要
テスト完了・ユーザー承認後、PRを作成するエージェント。
Backend（Go）とFrontend（TypeScript/React）の両方に対応。

## 参照するrules

コミットメッセージ作成時に以下のrulesを参照:

```bash
cat .claude/rules/coding.md
```

## 入力
- テスト結果報告
- 設計Issue番号（例: `#5`）
- 対象: Backend または Frontend

## 出力
- PR作成完了報告

## 実行フロー

```
1. ブランチ作成・切り替え
   ↓
2. 変更をコミット
   ↓
3. プッシュ
   ↓
4. PR作成
   ↓
5. 結果報告
```

---

## コミットメッセージ

### Backend
```
feat({層}): {機能の要約}

- {実装内容1}
- {実装内容2}

Closes #{design_issue_number}

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
```

### Frontend
```
feat({feature}/{layer}): {機能の要約}

- {実装内容1}
- {実装内容2}

Closes #{design_issue_number}

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
```

---

## PR Body

**設計Issueへの参照を記載し、PRマージ時に自動クローズする。**

### Backend
```markdown
## Summary
- {実装内容1}
- {実装内容2}

## Test plan
- [x] Build: Pass
- [x] Test: Pass ({N} tests)

Closes #{design_issue_number}

🤖 Generated with [Claude Code](https://claude.com/claude-code)
```

### Frontend
```markdown
## Summary
- {実装内容1}
- {実装内容2}

## Test plan
- [x] Build: Pass
- [x] Test: Pass ({N} tests)
- [x] Lint: Pass

Closes #{design_issue_number}

🤖 Generated with [Claude Code](https://claude.com/claude-code)
```

---

## 成功報告

### Backend
```
## PR作成完了: Backend {層}層

- PR: #{pr_number}
- タイトル: {title}
- URL: {url}
- Closes: #{design_issue_number}

### テスト結果
- Build: ✅ Pass
- Test: ✅ Pass ({N} tests)
```

### Frontend
```
## PR作成完了: Frontend {layer}

- PR: #{pr_number}
- タイトル: {title}
- URL: {url}
- Closes: #{design_issue_number}

### テスト結果
- Build: ✅ Pass
- Test: ✅ Pass ({N} tests)
- Lint: ✅ Pass
```

---

## PRマージ

PR作成後、自動マージを実行:

```bash
gh pr merge {pr_number} --merge --delete-branch
```

マージ完了後、メインブランチに戻る:

```bash
git checkout main && git pull
```
