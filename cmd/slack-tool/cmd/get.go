package cmd

import (
	"fmt"
	"os"

	"github.com/shellme/slack-tool/internal/config"
	"github.com/shellme/slack-tool/internal/slack"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [thread-url]",
	Short: "データ取得コマンド",
	Long:  "Slackからデータを取得するためのコマンドです。",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// 引数がある場合は直接スレッド取得処理を実行
			getThreadCmd.Run(cmd, args)
		} else {
			// 引数がない場合はヘルプを表示
			cmd.Help()
		}
	},
}

var getThreadCmd = &cobra.Command{
	Use:   "thread <slack-thread-url>",
	Short: "スレッドの内容を取得・整形",
	Long: `指定されたSlackスレッドのURLから会話内容を取得し、
AIへの入力に適した人間が読みやすいプレーンテキスト形式で整形して表示します。

例:
  slack-tool get thread "https://your-workspace.slack.com/archives/C12345678/p1234567890123456"
  slack-tool get thread "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --output thread.md
  slack-tool get thread "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --output thread.md --format markdown`,
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
		threadInfo, err := slack.ParseSlackURL(url)
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

		// スレッドの内容を取得
		messages, err := client.GetThreadReplies(threadInfo.ChannelID, threadInfo.Timestamp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: スレッドの取得に失敗しました: %v\n", err)
			os.Exit(1)
		}

		// フォーマッターを作成
		formatter := slack.NewFormatter(client)

		// メッセージを整形
		formatted, err := formatter.FormatThread(messages)
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
			fmt.Printf("スレッドの内容を %s に保存しました\n", outputFile)
		} else {
			// 結果を標準出力に表示
			fmt.Print(formatted)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getThreadCmd)

	// get コマンドのフラグ（省略形用）
	getCmd.Flags().StringP("output", "o", "", "出力ファイル名を指定（例: thread.md, thread.txt）。拡張子で形式を自動判定")
	getCmd.Flags().StringP("format", "f", "text", "出力形式を指定（text / markdown）。指定があれば拡張子より優先")

	// フラグの定義
	getThreadCmd.Flags().StringP("output", "o", "", "出力ファイル名を指定（例: thread.md, thread.txt）。拡張子で形式を自動判定")
	getThreadCmd.Flags().StringP("format", "f", "text", "出力形式を指定（text / markdown）。指定があれば拡張子より優先")
}
