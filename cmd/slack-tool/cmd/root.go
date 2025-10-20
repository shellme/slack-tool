package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// ビルド時に設定される変数
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "slack-tool",
	Short: "Slack操作を行うCLIツール",
	Long: `slack-toolは、Slackの様々な操作を行うCLIツールです。
AIへの入力に適した人間が読みやすいプレーンテキスト形式で整形したり、
Slackへの投稿やその他の操作をサポートします。

主な機能:
- Slackスレッドの内容を取得・整形
- Slack APIトークンの設定・管理
- 将来的に投稿機能などを追加予定

使用例:
  slack-tool get "https://your-workspace.slack.com/archives/C12345678/p1234567890123456"
  slack-tool post "Hello, world!" --channel C12345678
  slack-tool config set token "xoxp-xxxxxxxxxxxxxx-xxxxxxxx"
  slack-tool config show`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "エラーが発生しました: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// バージョン表示のカスタマイズ
	rootCmd.SetVersionTemplate(fmt.Sprintf("slack-tool version %s\ncommit: %s\nbuilt: %s\n", version, commit, date))
}
