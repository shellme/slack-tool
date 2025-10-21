package cmd

import (
	"fmt"
	"os"

	"github.com/shellme/slack-tool/internal/config"
	"github.com/shellme/slack-tool/internal/slack"
	slackgo "github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [message-url]",
	Short: "データ取得コマンド",
	Long:  "Slackからデータを取得するためのコマンドです。",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// 引数がある場合は直接メッセージ取得処理を実行
			getMessageCmd.Run(cmd, args)
		} else {
			// 引数がない場合はヘルプを表示
			cmd.Help()
		}
	},
}

var getMessageCmd = &cobra.Command{
	Use:   "message <slack-message-url>",
	Short: "メッセージの内容を取得・整形",
	Long: `指定されたSlackメッセージのURLから内容を取得し、
AIへの入力に適した人間が読みやすいプレーンテキスト形式で整形して表示します。

単一メッセージ、スレッド全体、スレッドの親メッセージのみなど、
様々な取得パターンに対応しています。

例:
  # 単一メッセージを取得
  slack-tool get message "https://your-workspace.slack.com/archives/C12345678/p1234567890123456"
  
  # スレッド全体を取得（--thread フラグ）
  slack-tool get message "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --thread
  
  # スレッドの親メッセージのみを取得（--parent フラグ）
  slack-tool get message "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --parent
  
  # ファイルに保存
  slack-tool get message "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --output message.md`,
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
		threadInfo, err := slack.ParseThreadURL(url)
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
		includeThread, _ := cmd.Flags().GetBool("thread")
		parentOnly, _ := cmd.Flags().GetBool("parent")

		var messages []slackgo.Message
		var getErr error

		if includeThread {
			// スレッド全体を取得
			// スレッド返信URLの場合は ThreadTimestamp を使用
			threadTimestamp := threadInfo.Timestamp
			if threadInfo.ThreadTimestamp != "" {
				threadTimestamp = threadInfo.ThreadTimestamp
			}
			messages, getErr = client.GetThreadReplies(threadInfo.ChannelID, threadTimestamp)
			if getErr != nil {
				fmt.Fprintf(os.Stderr, "エラー: スレッドの取得に失敗しました: %v\n", getErr)
				os.Exit(1)
			}
		} else if parentOnly {
			// スレッドの親メッセージのみを取得
			message, getErr := client.GetMessageInfo(threadInfo.ChannelID, threadInfo.Timestamp)
			if getErr != nil {
				fmt.Fprintf(os.Stderr, "エラー: メッセージの取得に失敗しました: %v\n", getErr)
				os.Exit(1)
			}
			messages = []slackgo.Message{*message}
		} else {
			// 単一メッセージを取得（デフォルト）
			message, getErr := client.GetMessageInfo(threadInfo.ChannelID, threadInfo.Timestamp)
			if getErr != nil {
				fmt.Fprintf(os.Stderr, "エラー: メッセージの取得に失敗しました: %v\n", getErr)
				os.Exit(1)
			}
			messages = []slackgo.Message{*message}
		}

		// フォーマッターを作成
		formatter := slack.NewFormatter(client)

		// メッセージを整形
		var formatted string
		if includeThread {
			formatted, err = formatter.FormatThread(messages)
		} else {
			formatted, err = formatter.FormatMessage(messages[0])
		}
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
			fmt.Printf("メッセージの内容を %s に保存しました\n", outputFile)
		} else {
			// 結果を標準出力に表示
			fmt.Print(formatted)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getMessageCmd)

	// get コマンドのフラグ（省略形用）
	getCmd.Flags().StringP("output", "o", "", "出力ファイル名を指定（例: message.md, message.txt）。拡張子で形式を自動判定")
	getCmd.Flags().StringP("format", "f", "text", "出力形式を指定（text / markdown）。指定があれば拡張子より優先")
	getCmd.Flags().BoolP("thread", "t", false, "スレッド全体を取得する（返信も含む）")
	getCmd.Flags().BoolP("parent", "p", false, "スレッドの親メッセージのみを取得する")

	// get message コマンドのフラグ
	getMessageCmd.Flags().StringP("output", "o", "", "出力ファイル名を指定（例: message.md, message.txt）。拡張子で形式を自動判定")
	getMessageCmd.Flags().StringP("format", "f", "text", "出力形式を指定（text / markdown）。指定があれば拡張子より優先")
	getMessageCmd.Flags().BoolP("thread", "t", false, "スレッド全体を取得する（返信も含む）")
	getMessageCmd.Flags().BoolP("parent", "p", false, "スレッドの親メッセージのみを取得する")
}
