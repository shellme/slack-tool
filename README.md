# slack-tool

Slackの様々な操作を行うCLIツールです。AIへの入力に適した人間が読みやすいプレーンテキスト形式で整形したり、Slackへの投稿やその他の操作をサポートします。

## 機能

- Slackスレッドの内容を取得・整形
- Slackチャンネルの内容を取得・整形（期間指定対応）
- Slackチャンネルにメッセージを投稿
- Slackスレッドに返信を投稿
- Slack投稿のリアクション一覧を取得
- Slack APIトークンの設定・管理

## クイックスタート

> [!IMPORTANT]
> 初回使用時は、Slack APIトークンの設定が必要です。詳細は[API設定方法](docs/API_SETUP.md)を参照してください。

### インストール

```bash
# tapを追加
brew tap shellme/slack-tool

# インストール
brew install slack-tool
```

### 初回セットアップ

```bash
# 1. インストール確認
slack-tool --version

# 2. Slack APIトークンを設定
slack-tool config set-token "xoxp-your-token-here"

# 3. 設定確認
slack-tool config show
```

### アップグレード

```bash
# slack-toolを最新版にアップグレード
brew upgrade slack-tool

# または、tapを更新してからアップグレード
brew update
brew upgrade slack-tool

# バージョン確認
slack-tool --version
```

### 基本的な使用方法

> [!TIP]
> 各コマンドには省略形があります。`slack-tool get "url"`のように短縮して使用できます。

```bash
# スレッドの内容を取得
slack-tool get "https://workspace.slack.com/archives/C12345678/p1234567890123456"

# チャンネルの内容を取得（直近100件）
slack-tool get channel "https://workspace.slack.com/archives/C12345678"

# メッセージを投稿
slack-tool post "こんにちは！" --channel "C12345678"

# リアクション一覧を取得
slack-tool reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456"
```

## 省略コマンド（便利な短縮形）

```bash
# スレッド取得の短縮形
slack-tool get "url"

# メッセージ投稿の短縮形
slack-tool post "メッセージ"

# リアクション取得の短縮形
slack-tool reactions "url"
```

## 詳しい使い方

[USAGE.md](docs/USAGE.md)から、利用可能なコマンドリストを参照してください。

## ドキュメント

- [詳細な使用方法](docs/USAGE.md) - 全コマンドの詳細説明
- [API設定方法](docs/API_SETUP.md) - Slack APIトークンの取得と設定
- [出力例とサンプル](docs/EXAMPLES.md) - 各コマンドの出力例
- [開発者向け情報](docs/DEVELOPMENT.md) - 開発・ビルド・リリース方法

## ライセンス

MIT License