package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/shellme/slack-tool/internal/config"
	"github.com/shellme/slack-tool/internal/slack"
	"github.com/spf13/cobra"
)

var postCmd = &cobra.Command{
	Use:   "post <message>",
	Short: "メッセージを投稿",
	Long: `指定されたSlackチャンネルにメッセージを投稿します。

例:
  slack-tool post "Hello, world!" --channel C12345678
  slack-tool post "This is a test message" --channel C12345678 --thread 1234567890.123456
  slack-tool post "Hey @john, can you review this?" --channel C12345678
  slack-tool post "スレッド返信です" --thread-url "https://workspace.slack.com/archives/C12345678/p1234567890123456"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		message := args[0]

		cm := config.NewConfigManager()
		cfg, err := cm.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: 設定の読み込みに失敗しました: %v\n", err)
			os.Exit(1)
		}

		if cfg.SlackToken == "" {
			fmt.Fprintf(os.Stderr, "エラー: Slack APIトークンが設定されていません。\n")
			fmt.Fprintf(os.Stderr, "以下のコマンドでトークンを設定してください:\n")
			fmt.Fprintf(os.Stderr, "  slack-tool config set token \"xoxp-xxxxxxxxxxxxxx-xxxxxxxx\"\n")
			os.Exit(1)
		}

		client := slack.NewClient(cfg.SlackToken)

		if err := client.TestConnection(); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: Slack APIへの接続に失敗しました: %v\n", err)
			os.Exit(1)
		}

		// チャンネルIDを取得
		channelID, _ := cmd.Flags().GetString("channel")
		threadURL, _ := cmd.Flags().GetString("thread-url")
		threadTimestamp, _ := cmd.Flags().GetString("thread")

		// スレッドURLが指定されている場合は新しいメソッドを使用
		if threadURL != "" {
			err := client.PostThreadReplyByURL(message, threadURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "エラー: スレッド返信の投稿に失敗しました: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("スレッド返信を投稿しました: %s\n", message)
			return
		} else if channelID == "" {
			fmt.Fprintf(os.Stderr, "エラー: チャンネルIDまたはスレッドURLが指定されていません。--channel または --thread-url フラグを使用してください。\n")
			os.Exit(1)
		}

		// チャンネルIDがURLの場合は解析
		if strings.HasPrefix(channelID, "https://") {
			channelInfo, err := slack.ParseChannelURL(channelID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "エラー: チャンネルURLの解析に失敗しました: %v\n", err)
				os.Exit(1)
			}
			channelID = channelInfo.ChannelID
		}

		// スレッド返信かどうかチェック
		if threadTimestamp != "" {
			// スレッド返信
			err := client.PostThreadReply(channelID, message, threadTimestamp)
			if err != nil {
				fmt.Fprintf(os.Stderr, "エラー: スレッド返信の投稿に失敗しました: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("スレッド返信を投稿しました: %s\n", message)
		} else {
			// 通常のメッセージ投稿
			err := client.PostMessage(channelID, message)
			if err != nil {
				fmt.Fprintf(os.Stderr, "エラー: メッセージの投稿に失敗しました: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("メッセージを投稿しました: %s\n", message)
		}
	},
}

func init() {
	rootCmd.AddCommand(postCmd)

	postCmd.Flags().StringP("channel", "c", "", "投稿先のチャンネルIDまたはURL")
	postCmd.Flags().StringP("thread", "t", "", "スレッド返信する場合のタイムスタンプ")
	postCmd.Flags().StringP("thread-url", "u", "", "スレッド返信する場合のスレッドURL")
}
