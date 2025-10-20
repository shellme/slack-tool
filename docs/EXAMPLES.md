# 出力例とサンプル

> [!NOTE]
> このドキュメントでは、各コマンドの実際の出力例を紹介します。コピー&ペーストしてテストにご活用ください。

## スレッド取得の出力例

### プレーンテキスト形式

```
--- Slackスレッドの内容 (2025/01/18 取得) ---

[2025-01-17 10:30:15][@alice]:
本番環境でユーザーAさんから「画像がアップロードできない」と報告がありました。
どなたか状況わかりますか？

[2025-01-17 10:32:01][@bob]:
ログを確認しましたが、S3バケットの権限エラーが出ていますね。
おそらく昨日のデプロイでIAMロールの設定がもれていそうです。

--- ここまで ---
```

### Markdown形式

```markdown
# Slackスレッドの内容 (2025/01/18 取得)

[2025-01-17 10:30:15][@alice]:
本番環境でユーザーAさんから「画像がアップロードできない」と報告がありました。
どなたか状況わかりますか？

[2025-01-17 10:32:01][@bob]:
ログを確認しましたが、S3バケットの権限エラーが出ていますね。
おそらく昨日のデプロイでIAMロールの設定がもれていそうです。
```

## チャンネル取得の出力例

### 通常のチャンネル取得

```
情報: 15件のメッセージを取得しました。
--- Slackチャンネルの内容 (2025/01/18 取得) ---
チャンネル: #general

[2025-01-18 09:00:00][@alice]:
おはようございます！

[2025-01-18 09:05:00][@bob]:
おはようございます！今日もよろしくお願いします。

--- ここまで ---
```

### 期間指定でのチャンネル取得

> [!TIP]
> 期間指定は複数の形式に対応しています。日付形式、日時形式、Unixタイムスタンプのいずれでも指定可能です。

```bash
# 2024年1月の投稿を取得
slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --oldest "2024-01-01" --latest "2024-01-31"

# 日時形式での指定
slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --oldest "2024-01-01T00:00:00" --latest "2024-01-31T23:59:59"

# Unixタイムスタンプでの指定
slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --oldest "1704067200" --latest "1735689599"

# 期間指定とファイル出力を組み合わせ
slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --oldest "2024-01-01" --latest "2024-01-31" --output january.md
```

## リアクション取得の出力例

### 通常形式（デフォルト）

```
:+1: (3人)
  - alice
  - bob
  - henry

:suteki-pink: (5人)
  - charlie
  - diana
  - eve
  - frank
  - grace
```

### シンプル形式（`--simple`）

```
:+1:
alice
bob
henry

:suteki-pink:
charlie
diana
eve
frank
grace
```

### メールアドレス形式

```
:+1: (3人)
  - alice@example.com
  - bob@example.com
  - henry@example.com

:suteki-pink: (5人)
  - charlie@example.com
  - diana@example.com
  - eve@example.com
  - frank@example.com
  - grace@example.com
```

### フィルタ機能

> [!TIP]
> フィルタ機能は特定のリアクションのみを取得したい場合に便利です。リアクション名は `:emoji:` の形式で指定してください。

```bash
# 特定のリアクションのみを取得
slack-tool reactions "url" --filter ":+1:"

# 出力例
:+1: (3人)
  - alice
  - bob
  - henry
```

## メッセージ投稿の例

### 基本的な投稿

```bash
# チャンネルにメッセージを投稿
slack-tool post "こんにちは！" --channel "C12345678"

# チャンネルURLで投稿
slack-tool post "こんにちは！" --channel "https://workspace.slack.com/archives/C12345678"
```

### スレッド返信

```bash
# タイムスタンプでスレッド返信
slack-tool post "返信です！" --thread "1234567890.123456"

# スレッドURLで返信
slack-tool post "返信です！" --thread-url "https://workspace.slack.com/archives/C12345678/p1234567890123456"
```

## ファイル出力の例

### スレッドをMarkdownファイルに保存

```bash
slack-tool get "url" --output thread.md --format markdown
```

### チャンネル履歴をテキストファイルに保存

```bash
slack-tool get channel "url" --output channel.txt --limit 100
```

### リアクション一覧をファイルに保存

```bash
slack-tool reactions "url" --output reactions.txt --simple
```

## エラー出力の例

### 無効なトークン

```
エラー: Slack APIへの接続に失敗しました: invalid_auth
```

### 権限不足

```
エラー: チャンネルの取得に失敗しました: missing_scope
```

### メッセージが見つからない

```
エラー: 指定されたタイムスタンプのメッセージが見つかりませんでした: 1234567890.123456
```
