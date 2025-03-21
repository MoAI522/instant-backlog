# スクラム開発バックログ管理システム

マークダウンファイルを使用してスクラム開発のバックログを管理するシンプルな CLI ツールです。

新しく参加する開発者は、[オンボーディングガイド](ONBOARDING.md)を参照してください。

## 機能

- マークダウンファイルで Epic と Issue を管理
- order.csv との自動同期による優先順位付け
- ステータス変更時のファイル名自動更新
- ファイル変更の自動監視と同期
- 関連する Issue がすべて Close になると Epic も自動的に Close に更新
- プロジェクトの初期化機能（テンプレートと使用方法ドキュメント付き）

## 使用方法

```bash
# order.csvを現在のオープンIssueと同期
./instant-backlog sync
# または省略形を使用
./ib sync

# Front Matterに基づいてファイル名を更新
./instant-backlog rename
# または省略形を使用
./ib rename

# ファイル変更を監視して自動的にsyncとrenameを実行
./instant-backlog watch [project_path]
# または省略形を使用
./ib watch [project_path]

# ファイル監視を停止
./instant-backlog unwatch [project_path]
# または省略形を使用
./ib unwatch [project_path]

# 新しいプロジェクトを初期化
./instant-backlog init [project_path]
# または省略形を使用
./ib init [project_path]
```

## ファイル構造

初期化コマンドでは以下の構造が自動的に作成されます：

```
projects/
  ├── epic/      # Epicファイル格納ディレクトリ
  ├── issues/    # Issueファイル格納ディレクトリ
  └── order.csv  # 実施順管理ファイル
```

## マークダウン形式

### Issue

```markdown
---
id: 1
title: "タスクのタイトル"
status: "Open" # "Open" または "Close"
epic: 3 # 関連するEpicのID
estimate: 5 # ポイント数
---

Issue 本文...
```

### Epic

```markdown
---
id: 3
title: "Epicのタイトル"
status: "Open" # "Open" または "Close"
---

Epic 本文...
```

## ビルド方法

```bash
go build -o dist/ib ./cmd/instant-backlog
```
