package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/shellme/slack-tool/internal/config"
	"github.com/shellme/slack-tool/internal/slack"
	"github.com/spf13/cobra"
)

var reactionsCmd = &cobra.Command{
	Use:   "get reactions <message-url>",
	Short: "指定した投稿のリアクション一覧を取得",
	Long: `指定したSlack投稿のリアクション一覧を取得します。

例:
  slack-tool get reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456"
  slack-tool get reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --filter ":参加します:"
  slack-tool get reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --email
  slack-tool get reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --output reactions.txt
  slack-tool get reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --filter "承知_しました" --email --output participants.txt`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		messageURL := args[0]

		cm := config.NewConfigManager()
		cfg, err := cm.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: 設定の読み込みに失敗しました: %v\n", err)
			os.Exit(1)
		}

		client := slack.NewClient(cfg.SlackToken)

		if err := client.TestConnection(); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: Slack APIへの接続に失敗しました: %v\n", err)
			os.Exit(1)
		}

		// オプションを取得
		filter, _ := cmd.Flags().GetString("filter")
		email, _ := cmd.Flags().GetBool("email")
		outputFile, _ := cmd.Flags().GetString("output")

		// リアクション一覧を取得
		reactions, err := client.GetReactions(messageURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: リアクションの取得に失敗しました: %v\n", err)
			os.Exit(1)
		}

		// フィルタが指定されている場合は適用
		if filter != "" {
			reactions = filterReactions(reactions, filter)
		}

		// 出力先を決定
		var output *os.File
		if outputFile != "" {
			file, err := os.Create(outputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "エラー: ファイルの作成に失敗しました: %v\n", err)
				os.Exit(1)
			}
			defer file.Close()
			output = file
		} else {
			output = os.Stdout
		}

		// ユーザー一覧を出力
		for _, reaction := range reactions {
			for _, user := range reaction.Users {
				if email {
					fmt.Fprintln(output, user.Email)
				} else {
					fmt.Fprintln(output, user.Name)
				}
			}
		}

		// ファイルに保存した場合のメッセージ
		if outputFile != "" {
			fmt.Fprintf(os.Stderr, "リアクション一覧を %s に保存しました\n", outputFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(reactionsCmd)

	reactionsCmd.Flags().StringP("filter", "f", "", "特定のリアクションのみをフィルタ（例: :参加します:）")
	reactionsCmd.Flags().BoolP("email", "e", false, "ユーザー名の代わりにメールアドレスを出力")
	reactionsCmd.Flags().StringP("output", "o", "", "結果をファイルに保存（例: reactions.txt）")
}

// filterReactions filters reactions by the specified emoji
func filterReactions(reactions []slack.ReactionInfo, filter string) []slack.ReactionInfo {
	// コロンが含まれている場合は除去
	cleanFilter := strings.Trim(filter, ":")

	var filtered []slack.ReactionInfo
	for _, reaction := range reactions {
		if reaction.Name == cleanFilter {
			filtered = append(filtered, reaction)
		}
	}
	return filtered
}
