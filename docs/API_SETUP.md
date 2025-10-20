# Slack API設定方法

> [!IMPORTANT]
> slack-toolを使用するには、Slack APIトークンの設定が必須です。以下の手順に従って設定してください。

## Slack APIトークンの取得

### 1. Slack APIにアクセス

1. [Slack API](https://api.slack.com/) にアクセス
2. 「Create an app」をクリック
3. 「From scratch」を選択
4. アプリ名とワークスペースを選択
5. 「OAuth & Permissions」に移動
6. 「User Token Scopes」にスコープを追加（後述）
7. 「Install to Workspace」をクリック
8. 生成された「User OAuth Token」をコピー

## 必要なスコープ（権限）

> [!WARNING]
> 以下のスコープがすべて必要です。権限不足の場合は適切に動作しません。

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

## トークンの設定

```bash
# トークンを設定
slack-tool config set-token "xoxp-your-token-here"

# 設定を確認
slack-tool config show
```

## エラーハンドリング

- 無効なトークンが設定されている場合、適切なエラーメッセージが表示されます
- 権限不足の場合は、必要なスコープの追加を促すメッセージが表示されます
- ネットワークエラーの場合は、接続状況を確認するよう案内されます
