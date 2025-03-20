# instant-backlog オンボーディングガイド

このドキュメントは、instant-backlog プロジェクトに新しく参加する開発者向けの包括的なガイドです。

## 1. プロジェクト概要

instant-backlog は、マークダウンファイルを使用してスクラム開発のバックログを管理するシンプルな CLI ツールです。以下の主要機能を提供します：

- マークダウンファイルで Epic と Issue を管理
- order.csv との自動同期による優先順位付け
- ステータス変更時のファイル名自動更新
- ファイル変更の自動監視と同期

## 2. プロジェクト構造

```
instant-backlog/
├── cmd/                    # アプリケーションエントリーポイント
│   └── instant-backlog/           # メインCLIツール
│       └── main.go
├── internal/               # 内部パッケージ（外部からインポート不可）
│   ├── commands/           # CLIコマンド実装
│   ├── config/             # 設定管理
│   ├── fileops/            # ファイル操作ユーティリティ
│   ├── models/             # データモデル
│   ├── parser/             # マークダウンパーサー
│   └── watcher/            # ファイル監視機能
├── pkg/                    # 外部パッケージ（外部からインポート可能）
│   └── utils/              # 汎用ユーティリティ
├── projects/               # プロジェクトデータディレクトリ
│   ├── epic/               # Epicファイル格納ディレクトリ
│   ├── issues/             # Issueファイル格納ディレクトリ
│   └── order.csv           # 実施順管理ファイル
└── test/                   # テストケース
```

## 3. 主要コンポーネント解説

### コマンドライン処理 (cmd/instant-backlog)

Cobra ライブラリを使用して CLI コマンドを実装しています。主要コマンド：

- `sync` - order.csv をオープン Issue と同期
- `rename` - Front Matter に基づいてファイル名を更新
- `watch` - ファイル変更を監視して自動で sync と rename を実行
- `unwatch` - ファイル監視を停止

### 内部パッケージ

- **config**: アプリケーション設定を管理。ディレクトリパスなどを含む
- **commands**: CLI コマンドの実装（`sync`、`rename`、`watch`、`unwatch`）
- **models**: Epic と issue のデータモデル
- **parser**: マークダウンファイルの Front Matter を解析
- **fileops**: ファイル操作のユーティリティ
- **watcher**: ファイル変更の監視と自動処理

### 外部パッケージ

- **utils**: ファイル名生成などの汎用ユーティリティ

### データ構造

#### Epic

```markdown
---
id: 1
title: "Epicのタイトル"
status: "Open" # "Open" または "Close"
---

Epic 本文...
```

#### Issue

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

#### order.csv

```csv
id,title,epic,estimate
1,最初のタスク,1,3
2,機能追加タスク,1,5
```

## 4. 開発環境セットアップ

### 前提条件

- Go 1.24 以上
- 任意のテキストエディタまたは IDE

### ローカル環境セットアップ

1. リポジトリのクローン

```bash
git clone <repository-url> instant-backlog
cd instant-backlog
```

2. 依存関係のインストール

```bash
go mod download
```

3. アプリケーションのビルド

```bash
go build -o instant-backlog ./cmd/instant-backlog
```

4. テストの実行

```bash
./run_tests.sh
```

## 5. デバッグとトラブルシューティング

- ファイル名の不一致エラー: `rename`コマンドを実行して状態とファイル名の同期を試みる
- order.csv の不整合: `sync`コマンドを実行して CSV ファイルと issue を同期

## 6. コントリビューションガイドライン

1. 新機能や修正は新しいブランチで開発
2. コミットメッセージは明確に変更内容を記述
3. プルリクエスト前にすべてのテストが通過していることを確認
4. コードは Go 標準のコーディング規約に従う

## 7. 主要な使用方法と例

### バックログの同期

```bash
./instant-backlog sync
# または省略形を使用
./ib sync
```

### ファイル名の更新

```bash
./instant-backlog rename
# または省略形を使用
./ib rename
```

### ファイル変更の監視

```bash
./instant-backlog watch [project_path]
# または省略形を使用
./ib watch [project_path]
```

### ファイル監視の停止

```bash
./instant-backlog unwatch [project_path]
# または省略形を使用
./ib unwatch [project_path]
```

## 8. よくある質問

**Q: ファイルのステータスはどのように変更しますか？**
A: マークダウンファイルの Front Matter で`status`フィールドを"Open"または"Close"に変更し、`rename`コマンドを実行します。

**Q: epic と issue の関連付けはどのように行いますか？**
A: issue の Front Matter で`epic`フィールドに関連する epic の ID を指定します。

**Q: 自動ファイル監視はどのように設定しますか？**
A: `watch`コマンドを実行してプロジェクトディレクトリを監視します。例: `./ib watch projects`。監視を停止するには Ctrl+C を押すか、別のターミナルで`unwatch`コマンドを実行します。
