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

### Homebrew（推奨）

```bash
# tapを追加
brew tap shellme/slack-tool

# インストール
brew install slack-tool
```

### Go install

```bash
go install github.com/shellme/slack-tool@latest
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

### スレッドの取得

SlackスレッドのURLを指定して、会話内容を取得・整形します。

```bash
# 標準出力に表示
slack-tool get "https://your-workspace.slack.com/archives/C12345678/p1234567890123456"

# Markdownファイルとして保存
slack-tool get "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --output thread.md

# プレーンテキスト形式で保存
slack-tool get "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --output thread.txt --format text

# Markdown形式で保存（明示的に指定）
slack-tool get "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --output thread.md --format markdown
```

### チャンネルの取得

SlackチャンネルのURLを指定して、会話内容を取得・整形します。

```bash
# 標準出力に表示（最新100件）
slack-tool channel "https://your-workspace.slack.com/archives/C12345678"

# 最新50件を取得
slack-tool channel "https://your-workspace.slack.com/archives/C12345678" --limit 50

# Markdownファイルとして保存
slack-tool channel "https://your-workspace.slack.com/archives/C12345678" --output channel.md

# プレーンテキスト形式で保存
slack-tool channel "https://your-workspace.slack.com/archives/C12345678" --output channel.txt --format text

# プライベートチャンネルやダイレクトメッセージも取得可能
slack-tool channel "https://your-workspace.slack.com/archives/G12345678" --limit 200
```

## 出力例

### プレーンテキスト形式
```
--- Slackスレッドの内容 (2025/01/18 取得) ---

[2025-01-17 10:30:15][@suzuki_san]:
本番環境でユーザーAさんから「画像がアップロードできない」と報告がありました。
どなたか状況わかりますか？

[2025-01-17 10:32:01][@sato_san]:
ログを確認しましたが、S3バケットの権限エラーが出ていますね。
おそらく昨日のデプロイでIAMロールの設定がもれていそうです。

--- ここまで ---
```

### Markdown形式
```markdown
# Slackスレッドの内容 (2025/01/18 取得)

[2025-01-17 10:30:15][@suzuki_san]:
本番環境でユーザーAさんから「画像がアップロードできない」と報告がありました。
どなたか状況わかりますか？

[2025-01-17 10:32:01][@sato_san]:
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
   - `reactions:read`
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
- `slack-tool get <slack-thread-url>` - Slackスレッドの内容を取得・整形
- `slack-tool get <slack-thread-url> --output <file>` - スレッドの内容をファイルに保存
- `slack-tool get <slack-thread-url> --output <file> --format markdown` - Markdown形式で保存
- `slack-tool channel <slack-channel-url>` - Slackチャンネルの内容を取得・整形
- `slack-tool channel <slack-channel-url> --output <file>` - チャンネルの内容をファイルに保存
- `slack-tool channel <slack-channel-url> --limit <number>` - 取得件数を指定

### メッセージ投稿コマンド
- `slack-tool post "<message>" --channel <channel-id>` - チャンネルにメッセージを投稿
- `slack-tool post "<message>" --channel <channel-url>` - チャンネルURLでメッセージを投稿
- `slack-tool post "<message>" --thread-url <thread-url>` - スレッドに返信を投稿

### リアクション管理コマンド
- `slack-tool reactions <message-url>` - 指定した投稿のリアクション一覧を取得
- `slack-tool reactions <message-url> --filter ":emoji:"` - 特定のリアクションのみをフィルタ
- `slack-tool reactions <message-url> --email` - メールアドレス形式で出力

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

