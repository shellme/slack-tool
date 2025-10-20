package cmd

import (
	"fmt"
	"os"

	"github.com/shellme/slack-tool/internal/config"
	"github.com/shellme/slack-tool/internal/slack"
	"github.com/spf13/cobra"
)

var channelCmd = &cobra.Command{
	Use:   "get channel <slack-channel-url>",
	Short: "チャンネルの内容を取得・整形",
	Long: `指定されたSlackチャンネルのURLから会話内容を取得し、
AIへの入力に適した人間が読みやすいプレーンテキスト形式で整形して表示します。

例:
  slack-tool get channel "https://your-workspace.slack.com/archives/C12345678"
  slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --output channel.md
  slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --output channel.md --format markdown
  slack-tool get channel "https://your-workspace.slack.com/archives/C12345678" --limit 50`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		// 設定を読み込み
		cm := config.NewConfigManager()
		cfg, err := cm.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: 設定の読み込みに失敗しました: %v\n", err)
			os.Exit(1)
		}

		// トークンが設定されているかチェック
		if cfg.SlackToken == "" {
			fmt.Fprintf(os.Stderr, "エラー: Slack APIトークンが設定されていません。\n")
			fmt.Fprintf(os.Stderr, "以下のコマンドでトークンを設定してください:\n")
			fmt.Fprintf(os.Stderr, "  slack-tool config set token \"xoxp-xxxxxxxxxxxxxx-xxxxxxxx\"\n")
			os.Exit(1)
		}

		// URLを解析
		channelInfo, err := slack.ParseChannelURL(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}

		// Slackクライアントを作成
		client := slack.NewClient(cfg.SlackToken)

		// 接続をテスト
		if err := client.TestConnection(); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: Slack APIへの接続に失敗しました: %v\n", err)
			os.Exit(1)
		}

		// 取得件数を取得
		limit, _ := cmd.Flags().GetInt("limit")

		// チャンネルの内容を取得（スレッド返信も含む）
		messages, err := client.GetChannelHistoryWithThreads(channelInfo.ChannelID, limit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: チャンネルの取得に失敗しました: %v\n", err)
			os.Exit(1)
		}

		// チャンネル情報を取得
		channel, err := client.GetChannelInfo(channelInfo.ChannelID)
		var channelName string
		if err != nil {
			// チャンネル情報が取得できない場合はIDをそのまま使用
			channelName = channelInfo.ChannelID
		} else {
			channelName = channel.Name
		}

		// フォーマッターを作成
		formatter := slack.NewFormatter(client)

		// メッセージを整形
		formatted, err := formatter.FormatChannel(messages, channelName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: メッセージの整形に失敗しました: %v\n", err)
			os.Exit(1)
		}

		// 出力ファイルが指定されているかチェック
		outputFile, _ := cmd.Flags().GetString("output")
		format, _ := cmd.Flags().GetString("format")

		if outputFile != "" {
			// ファイルに保存
			err := saveToFile(formatted, outputFile, format)
			if err != nil {
				fmt.Fprintf(os.Stderr, "エラー: ファイルの保存に失敗しました: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("チャンネルの内容を %s に保存しました\n", outputFile)
		} else {
			// 結果を標準出力に表示
			fmt.Print(formatted)
		}
	},
}

func init() {
	rootCmd.AddCommand(channelCmd)

	// フラグの定義
	channelCmd.Flags().StringP("output", "o", "", "出力ファイル名を指定（例: channel.md, channel.txt）。拡張子で形式を自動判定")
	channelCmd.Flags().StringP("format", "f", "text", "出力形式を指定（text / markdown）。指定があれば拡張子より優先")
	channelCmd.Flags().IntP("limit", "l", 100, "取得するメッセージ数を指定（デフォルト: 100）")
}
