# スプリントバックログ運用ガイド

このドキュメントでは、Instant Backlogを使用したスプリントバックログの運用方法について説明します。

## 1. プロジェクト構成

Instant Backlogの各プロジェクトは以下の構成になっています：

```
projects/
  ├── epic/      # Epicファイル格納ディレクトリ
  ├── issues/    # Issueファイル格納ディレクトリ
  ├── order.csv  # 実施順管理ファイル
  └── README.md  # このファイル (プロジェクト説明)
```

## 2. ファイルフォーマット

### Epic

Epicファイルは以下のフォーマットで作成します：

```markdown
---
id: 1
title: "エピックのタイトル"
status: "Open"  # "Open" または "Close"
---

エピックの詳細な説明...
```

### Issue

Issueファイルは以下のフォーマットで作成します：

```markdown
---
id: 1
title: "課題のタイトル"
status: "Open"  # "Open" または "Close"
epic: 1         # 関連するエピックのID
estimate: 3     # ストーリーポイント（工数見積もり）
---

課題の詳細な説明...
```

### order.csv

order.csvファイルは優先順位付けのためのファイルで、以下のフォーマットで作成します：

```csv
id,title,epic,estimate
1,最初の課題,1,3
2,次の課題,1,5
```

- **id**: IssueのID (Epicは含まれません)
- **title**: Issueのタイトル
- **epic**: 関連するEpicのID
- **estimate**: ストーリーポイント（工数見積もり）

**重要**: order.csvファイルは上から順に実行優先度が高いことを意味します。

## 3. ファイル命名規則

ファイル名は以下の規則に従って自動的に生成されます：

```
{ID}_{STATUS_PREFIX}_{TITLE}.md
```

- **ID**: タスクまたはEpicの一意のID
- **STATUS_PREFIX**: ステータスの頭文字（`O` = Open、`C` = Close）
- **TITLE**: タイトル（スペースはアンダースコアに置換、特殊文字は除去）

例:
- `1_O_最初のエピック.md` (ID:1、ステータス:Open、タイトル:最初のエピック)
- `2_C_完了したタスク.md` (ID:2、ステータス:Close、タイトル:完了したタスク)

## 4. 運用方法

### 4.1 新しいEpicの作成

1. epic/ディレクトリに新しいファイルを作成
2. 上記のEpicフォーマットに従ってファイルを記入
3. `rename`コマンドを実行（または自動監視中であれば自動的に実行）

### 4.2 新しいIssueの作成

1. issues/ディレクトリに新しいファイルを作成
2. 上記のIssueフォーマットに従ってファイルを記入
3. `rename`コマンドを実行
4. `sync`コマンドを実行してorder.csvに追加

### 4.3 Issueのステータス変更

1. 該当するIssueのFront Matterの`status`フィールドを変更
2. `rename`コマンドを実行（ファイル名が自動更新）
3. ステータスを"Close"にした場合、`sync`コマンドを実行（order.csvから削除）

### 4.4 Epicのステータス変更

1. 該当するEpicのFront Matterの`status`フィールドを変更
2. `rename`コマンドを実行（ファイル名が自動更新）

**注**: 関連するすべてのIssueがCloseになると、Epicも自動的にCloseになります。

### 4.5 優先順位の変更

1. order.csvファイルを編集して優先順位を変更（上にあるほど優先度が高い）
2. 手動で編集するか、`sync`コマンドの実行後に並べ替えます

## 5. 自動監視モード

プロジェクトの変更を自動的に監視して同期するには：

```bash
./instant-backlog watch projects
# または省略形を使用
./ib watch projects
```

監視を停止するには:

```bash
./instant-backlog unwatch projects
# または省略形を使用
./ib unwatch projects
```

## 6. よくある質問

**Q: 同じIDのEpicやIssueが複数存在する場合どうなりますか？**
A: 同じIDを持つ複数のファイルが存在する場合、Close状態のファイルが優先されます。同期コマンド実行時にチェックが行われ、警告メッセージが表示されます。

**Q: ファイル名を手動で変更してもいいですか？**
A: 避けてください。Front Matterの内容を編集し、`rename`コマンドを実行することをお勧めします。手動変更すると不整合が発生する可能性があります。

**Q: Epicを削除するとどうなりますか？**
A: Epicを削除しても関連するIssueは削除されません。関連するIssueのepicフィールドを更新する必要があります。

**Q: order.csvにないIssueはどうなりますか？**
A: `sync`コマンドを実行すると、OpenステータスのIssueはすべてorder.csvに追加され、CloseステータスのIssueは削除されます。
