# slack-tool

Slackの様々な操作を行うCLIツールです。AIへの入力に適した人間が読みやすいプレーンテキスト形式で整形したり、Slackへの投稿やその他の操作をサポートします。

## 機能

- Slackスレッドの内容を取得・整形
- Slackチャンネルの内容を取得・整形
- Slackチャンネルにメッセージを投稿
- Slackスレッドに返信を投稿
- Slack投稿のリアクション一覧を取得
- Slack APIトークンの設定・管理

## インストール

```bash
# tapを追加
brew tap shellme/slack-tool

# インストール
brew install slack-tool
```

## アップデート

```bash
# Homebrewを更新
brew update

# slack-toolをアップデート
brew upgrade slack-tool

# または全パッケージをアップデート
brew upgrade
```

## 初回セットアップ

### 1. インストール確認
```bash
# インストールが成功しているか確認
slack-tool --version

# 期待される出力例
# slack-tool version v0.1.1
```

### 2. Slack APIトークンの取得
Slack APIのUser Tokenを取得します（詳細は後述の「Slack APIトークンの取得方法」を参照）。

### 3. トークンの設定
取得したトークンを設定します：

```bash
slack-tool config set token "xoxp-xxxxxxxxxxxxxx-xxxxxxxx"
```

### 4. 設定の確認
```bash
slack-tool config show
```

## 使用方法

### スレッドの取得（get thread）

SlackスレッドのURLを指定して、会話内容を取得・整形します。

```bash
# 標準出力に表示（省略形）
slack-tool get "https://your-workspace.slack.com/archives/C12345678/p1234567890123456"

# 標準出力に表示（完全形）
slack-tool get thread "https://your-workspace.slack.com/archives/C12345678/p1234567890123456"

# Markdownファイルとして保存
slack-tool get thread "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --output thread.md

# プレーンテキスト形式で保存
slack-tool get thread "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --output thread.txt

# Markdown形式で保存（明示的に指定）
slack-tool get thread "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --output thread.md --format markdown
```

補足:
- **出力形式の自動判定**: `--output` の拡張子で形式を自動判定します（`.md`/`.markdown` → markdown、それ以外 → text）。
- **明示的指定の優先**: `--format` を指定した場合は拡張子より `--format` が優先されます。

### チャンネルの取得（get channel）

SlackチャンネルのURLを指定して、会話内容を取得・整形します。

```bash
# 標準出力に表示（最新100件）
slack-tool get channel "https://your-workspace.slack.com/archives/C12345678"

# 最新50件を取得
slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --limit 50

# Markdownファイルとして保存
slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --output channel.md

# プレーンテキスト形式で保存
slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --output channel.txt

# プライベートチャンネルやダイレクトメッセージも取得可能
slack-tool get channel "https://your-workspace.slack.com/archives/G12345678" --limit 200

# 2024年1月の投稿を取得
slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --oldest "2024-01-01" --latest "2024-01-31"

# 日時形式での指定
slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --oldest "2024-01-01T00:00:00" --latest "2024-01-31T23:59:59"

# Unixタイムスタンプでの指定
slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --oldest "1704067200" --latest "1735689599"

# 期間指定とファイル出力を組み合わせ
slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --oldest "2024-01-01" --latest "2024-01-31" --output january.md
```

### メッセージの投稿（post message）

Slackチャンネルにメッセージを投稿します。

```bash
# 省略形でメッセージを投稿
slack-tool post "Hello, world!" --channel C12345678

# 完全形でメッセージを投稿
slack-tool post message "Hello, world!" --channel C12345678

# スレッドに返信
slack-tool post "スレッド返信です" --thread-url "https://workspace.slack.com/archives/C12345678/p1234567890123456"

# チャンネルURLでメッセージを投稿
slack-tool post "This is a test message" --channel "https://your-workspace.slack.com/archives/C12345678"
```

### リアクションの取得（get reactions）

指定した投稿のリアクション一覧を取得します。

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

#### 出力例

**通常形式（デフォルト）:**
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

**シンプル形式（`--simple`）:**
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

補足:
- **出力形式の自動判定**: `--output` の拡張子で形式を自動判定します（`.md`/`.markdown` → markdown、それ以外 → text）。
- **明示的指定の優先**: `--format` を指定した場合は拡張子より `--format` が優先されます。
- **取得件数制限**: 1回のリクエストで最大1,000件まで取得可能です。それ以上の取得が必要な場合は期間指定（`--oldest`/`--latest`）を使用して複数回に分けて取得してください。
- **リアクション統合**: スキントーンなどの修飾子（`:skin-tone-2:`など）は基本のリアクション名に統合されます。例：`:+1:` と `:+1::skin-tone-2:` は `:+1:` として集計されます。

## 出力例

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

## Slack APIトークンの取得方法

1. [Slack API](https://api.slack.com/) にアクセス
2. 「Create an app」をクリック
3. 「From scratch」を選択
4. アプリ名とワークスペースを選択
5. 「OAuth & Permissions」に移動
6. 「User Token Scopes」にスコープを追加（後述）
7. 「Install to Workspace」をクリック
8. 生成された「User OAuth Token」をコピー

## 必要なスコープ（権限）

- `channels:history` - パブリックチャンネルの履歴を読み取り
- `groups:history` - プライベートチャンネルの履歴を読み取り
- `im:history` - ダイレクトメッセージの履歴を読み取り
- `mpim:history` - マルチパーティダイレクトメッセージの履歴を読み取り
- `users:read` - ユーザー情報を読み取り
- `usergroups:read` - ユーザーグループ情報を読み取り
- `reactions:read` - リアクション情報を読み取り
- `chat:write` - メッセージを投稿
- `chat:write.public` - パブリックチャンネルにメッセージを投稿
- `chat:write.customize` - メッセージのカスタマイズ
- `reactions:read` - リアクションを読み取り

## 設定ファイル

設定は `~/.config/slack-tool/config.json` に保存されます。

```json
{
  "slack_token": "xoxp-xxxxxxxxxxxxxx-xxxxxxxx"
}
```

## エラーハンドリング

- 無効なURLが渡された場合、適切なエラーメッセージを表示
- APIトークンが未設定または無効な場合、エラーメッセージを表示
- ネットワークエラーやAPIレート制限にも対応

## コマンド一覧

### 基本コマンド
- `slack-tool --help` - ヘルプを表示
- `slack-tool --version` - バージョンを表示

### 設定コマンド
- `slack-tool config set token <token>` - Slack APIトークンを設定
- `slack-tool config show` - 現在の設定を表示

### データ取得コマンド
- `slack-tool get thread <slack-thread-url>` - Slackスレッドの内容を取得・整形
- `slack-tool get thread <slack-thread-url> --output <file>` - スレッドの内容をファイルに保存
- `slack-tool get thread <slack-thread-url> --output <file> --format markdown` - Markdown形式で保存
- `slack-tool get channel <slack-channel-url>` - Slackチャンネルの内容を取得・整形
- `slack-tool get channel <slack-channel-url> --output <file>` - チャンネルの内容をファイルに保存
- `slack-tool get channel <slack-channel-url> --limit <number>` - 取得件数を指定
- `slack-tool get channel <slack-channel-url> --oldest <date> --latest <date>` - 期間を指定して取得

### メッセージ投稿コマンド（post）
- `slack-tool post message "<message>" --channel <channel-id>` - チャンネルにメッセージを投稿
- `slack-tool post message "<message>" --channel <channel-url>` - チャンネルURLでメッセージを投稿
- `slack-tool post message "<message>" --thread-url <thread-url>` - スレッドに返信を投稿
- `slack-tool post message "<message>" --channel <channel-id> --thread <timestamp>` - TSを直接指定してスレッド返信

### リアクション管理コマンド（get reactions）
- `slack-tool get reactions <message-url>` - 指定した投稿のリアクション一覧を取得
- `slack-tool get reactions <message-url> --filter ":emoji:"` - 特定のリアクションのみをフィルタ
- `slack-tool get reactions <message-url> --email` - メールアドレス形式で出力
- `slack-tool get reactions <message-url> --output <file>` - 結果をファイルに保存

### 省略コマンド（便利な短縮形）
- `slack-tool get <slack-thread-url>` - `get thread` の省略形
- `slack-tool post "<message>" --channel <channel-id>` - `post message` の省略形
- `slack-tool reactions <message-url>` - `get reactions` の省略形

## 開発

### 前提条件
- Go 1.22以上
- [github.com/spf13/cobra](https://github.com/spf13/cobra) - CLIフレームワーク
- [github.com/slack-go/slack](https://github.com/slack-go/slack) - Slack公式Goクライアント

### ビルド
```bash
# ビルド
make build

# インストール
make install

# テスト
make test
```

### 開発ワークフロー

#### 新機能の追加
1. 機能ブランチを作成
2. 機能を実装
3. テストを追加
4. ドキュメントを更新
5. プルリクエストを作成

#### リリース
1. バージョンタグを作成
2. GitHub Actionsが自動的にリリースをビルド・公開
3. Homebrew Formulaが自動更新

詳細なリリース手順については、[リリースワークフロー](docs/RELEASE_WORKFLOW.md)を参照してください。

## ライセンス

MIT License

