# 使用方法

> [!NOTE]
> このドキュメントでは、slack-toolの全コマンドの詳細な使用方法を説明します。基本的な使用方法は[メインREADME](../README.md)を参照してください。

## コマンド一覧

### 基本コマンド

- `slack-tool --help` - ヘルプを表示
- `slack-tool --version` - バージョン情報を表示

### 設定コマンド

- `slack-tool config set token <token>` - Slack APIトークンを設定
- `slack-tool config show` - 現在の設定を表示

### データ取得コマンド（get）

#### メッセージの取得（get message）

> [!TIP]
> 最も汎用的なコマンドです。単一メッセージ、スレッド全体、スレッドの親メッセージのみなど、様々な取得パターンに対応しています。

指定したSlackメッセージのURLから内容を取得し、人間が読みやすくAIへの入力にも利用できる形式に整形して表示します。

```bash
# 省略形でメッセージ取得（単一メッセージ）
slack-tool get "https://workspace.slack.com/archives/C12345678/p1234567890123456"

# 完全形でメッセージ取得（単一メッセージ）
slack-tool get message "https://workspace.slack.com/archives/C12345678/p1234567890123456"

# スレッド全体を取得（返信も含む）
slack-tool get message "https://workspace.slack.com/archives/C12345678/p1234567890123456" --thread

# スレッドの親メッセージのみを取得
slack-tool get message "https://workspace.slack.com/archives/C12345678/p1234567890123456" --parent

# ファイルに保存
slack-tool get message "https://workspace.slack.com/archives/C12345678/p1234567890123456" --output message.md

# Markdown形式で保存
slack-tool get message "https://workspace.slack.com/archives/C12345678/p1234567890123456" --format markdown --output message.md
```


#### チャンネルの取得（get channel）

> [!IMPORTANT]
> 1回のリクエストで最大1,000件まで取得可能です。それ以上の取得が必要な場合は期間指定（`--oldest`/`--latest`）を使用して複数回に分けて取得してください。

指定したSlackチャンネルの内容を取得し、整形して表示します。期間指定や取得件数の制限も可能です。

```bash
# 省略形でチャンネル取得
slack-tool channel "https://workspace.slack.com/archives/C12345678"

# 完全形でチャンネル取得
slack-tool get channel "https://workspace.slack.com/archives/C12345678"

# 取得件数を指定
slack-tool get channel "https://workspace.slack.com/archives/C12345678" --limit 50

# 期間を指定して取得
slack-tool get channel "https://workspace.slack.com/archives/C12345678" --oldest "2024-01-01" --latest "2024-01-31"

# ファイルに保存
slack-tool get channel "https://workspace.slack.com/archives/C12345678" --output channel.md
```

#### リアクションの取得（get reactions）

> [!TIP]
> リアクションは自動的に統合されます。`:+1:` と `:+1::skin-tone-2:` は同じリアクションとして集計されます。Googleカレンダーなどにコピーする場合は `--simple` フラグが便利です。

指定した投稿のリアクション一覧を取得します。リアクション別にユーザーを整理して表示し、スキントーンなどの修飾子も統合されます。

```bash
# 省略形でリアクション一覧を取得
slack-tool reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456"

# 完全形でリアクション一覧を取得
slack-tool get reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456"

# 特定のリアクションのみをフィルタ
slack-tool reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --filter ":参加します:"

# メールアドレス形式で出力
slack-tool reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --email

# 結果をファイルに保存
slack-tool reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --output reactions.txt

# シンプル形式で出力（Googleカレンダーなどにコピーしやすい）
slack-tool reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --simple

# シンプル形式でメールアドレス出力
slack-tool reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --simple --email
```

### メッセージ投稿コマンド（post）

#### メッセージの投稿（post message）

> [!NOTE]
> チャンネルIDまたはチャンネルURLのどちらでも指定できます。スレッド返信の場合は、タイムスタンプまたはスレッドURLを指定してください。

指定したチャンネルにメッセージを投稿します。スレッド返信も可能です。

```bash
# 省略形でメッセージ投稿
slack-tool post "こんにちは！" --channel "C12345678"

# 完全形でメッセージ投稿
slack-tool post message "こんにちは！" --channel "C12345678"

# チャンネルURLで投稿
slack-tool post "こんにちは！" --channel "https://workspace.slack.com/archives/C12345678"

# スレッドに返信
slack-tool post "返信です！" --thread "1234567890.123456"

# スレッドURLで返信
slack-tool post "返信です！" --thread-url "https://workspace.slack.com/archives/C12345678/p1234567890123456"
```

## フラグ一覧

### 共通フラグ

- `--output`, `-o` - 出力ファイル名を指定
- `--format`, `-f` - 出力形式を指定（text / markdown）

### get message 専用フラグ

- `--thread`, `-t` - スレッド全体を取得する（返信も含む）
- `--parent`, `-p` - スレッドの親メッセージのみを取得する

### get channel 専用フラグ

- `--limit`, `-l` - 取得するメッセージ数を指定（デフォルト: 100）
- `--oldest` - 取得開始日時を指定
- `--latest` - 取得終了日時を指定

### get reactions 専用フラグ

- `--filter`, `-f` - 特定のリアクションのみをフィルタ
- `--email`, `-e` - ユーザー名の代わりにメールアドレスを出力
- `--simple`, `-s` - シンプル形式で出力（改行のみで区切り）

### post message 専用フラグ

- `--channel`, `-c` - 投稿先のチャンネルIDまたはURL
- `--thread`, `-t` - スレッド返信する場合のタイムスタンプ
- `--thread-url`, `-u` - スレッド返信する場合のスレッドURL

## 補足情報

- **出力形式の自動判定**: `--output` の拡張子で形式を自動判定します（`.md`/`.markdown` → markdown、それ以外 → text）。
- **明示的指定の優先**: `--format` を指定した場合は拡張子より `--format` が優先されます。
- **取得件数制限**: 1回のリクエストで最大1,000件まで取得可能です。それ以上の取得が必要な場合は期間指定（`--oldest`/`--latest`）を使用して複数回に分けて取得してください。
- **リアクション統合**: スキントーンなどの修飾子（`:skin-tone-2:`など）は基本のリアクション名に統合されます。例：`:+1:` と `:+1::skin-tone-2:` は `:+1:` として集計されます。
