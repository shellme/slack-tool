package cmd

import (
	"fmt"
	"os"

	"github.com/shellme/slack-tool/internal/config"
	"github.com/shellme/slack-tool/internal/slack"
	"github.com/spf13/cobra"
)

var channelCmd = &cobra.Command{
	Use:   "channel [channel-url]",
	Short: "チャンネルの内容を取得・整形",
	Long:  "Slackチャンネルの内容を取得するためのコマンドです。",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// 引数がある場合は直接チャンネル取得処理を実行
			getChannelCmd.Run(cmd, args)
		} else {
			// 引数がない場合はヘルプを表示
			cmd.Help()
		}
	},
}

var getChannelCmd = &cobra.Command{
	Use:   "channel <slack-channel-url>",
	Short: "チャンネルの内容を取得・整形",
	Long: `指定されたSlackチャンネルのURLから会話内容を取得し、
AIへの入力に適した人間が読みやすいプレーンテキスト形式で整形して表示します。

例:
  slack-tool channel "https://your-workspace.slack.com/archives/C12345678"
  slack-tool get channel "https://your-workspace.slack.com/archives/C12345678"
  slack-tool channel "https://your-workspace.slack.com/archives/C12345678" --output channel.md
  slack-tool channel "https://your-workspace.slack.com/archives/C12345678" --output channel.md --format markdown
  slack-tool channel "https://your-workspace.slack.com/archives/C12345678" --limit 50`,
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

		// フラグを取得
		limit, _ := cmd.Flags().GetInt("limit")
		oldest, _ := cmd.Flags().GetString("oldest")
		latest, _ := cmd.Flags().GetString("latest")

		// 1000件を超える場合の警告表示
		if limit > 1000 {
			fmt.Fprintf(os.Stderr, "警告: 指定された取得件数（%d件）はSlack APIの制限（1,000件）を超えています。\n", limit)
			fmt.Fprintf(os.Stderr, "実際には1,000件までしか取得されません。\n")
			fmt.Fprintf(os.Stderr, "大量のデータを取得する場合は期間指定（--oldest/--latest）を使用して複数回に分けて取得することをお勧めします。\n\n")
		}

		// チャンネルの内容を取得（スレッド返信も含む）
		messages, err := client.GetChannelHistoryWithThreadsInRange(channelInfo.ChannelID, limit, oldest, latest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: チャンネルの取得に失敗しました: %v\n", err)
			os.Exit(1)
		}

		// 取得件数の情報表示
		actualCount := len(messages)
		fmt.Fprintf(os.Stderr, "情報: %d件のメッセージを取得しました。\n", actualCount)

		if actualCount >= 1000 {
			fmt.Fprintf(os.Stderr, "情報: 1,000件の制限に達しました。古いメッセージが取得されていない可能性があります。\n")
			fmt.Fprintf(os.Stderr, "全期間のデータが必要な場合は期間指定（--oldest/--latest）を使用して複数回に分けて取得してください。\n\n")
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
	getCmd.AddCommand(getChannelCmd)

	// channel コマンドのフラグ（省略形用）
	channelCmd.Flags().StringP("output", "o", "", "出力ファイル名を指定（例: channel.md, channel.txt）。拡張子で形式を自動判定")
	channelCmd.Flags().StringP("format", "f", "text", "出力形式を指定（text / markdown）。指定があれば拡張子より優先")
	channelCmd.Flags().IntP("limit", "l", 100, "取得するメッセージ数を指定（デフォルト: 100）")
	channelCmd.Flags().StringP("oldest", "", "", "取得開始日時を指定（例: 2024-01-01, 2024-01-01T00:00:00, 1704067200）")
	channelCmd.Flags().StringP("latest", "", "", "取得終了日時を指定（例: 2024-12-31, 2024-12-31T23:59:59, 1735689599）")

	// get channel コマンドのフラグ
	getChannelCmd.Flags().StringP("output", "o", "", "出力ファイル名を指定（例: channel.md, channel.txt）。拡張子で形式を自動判定")
	getChannelCmd.Flags().StringP("format", "f", "text", "出力形式を指定（text / markdown）。指定があれば拡張子より優先")
	getChannelCmd.Flags().IntP("limit", "l", 100, "取得するメッセージ数を指定（デフォルト: 100）")
	getChannelCmd.Flags().StringP("oldest", "", "", "取得開始日時を指定（例: 2024-01-01, 2024-01-01T00:00:00, 1704067200）")
	getChannelCmd.Flags().StringP("latest", "", "", "取得終了日時を指定（例: 2024-12-31, 2024-12-31T23:59:59, 1735689599）")
}
