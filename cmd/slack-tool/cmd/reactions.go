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
	Use:   "reactions [message-url]",
	Short: "指定した投稿のリアクション一覧を取得",
	Long:  "Slack投稿のリアクション一覧を取得するためのコマンドです。",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// 引数がある場合は直接リアクション取得処理を実行
			getReactionsCmd.Run(cmd, args)
		} else {
			// 引数がない場合はヘルプを表示
			cmd.Help()
		}
	},
}

var getReactionsCmd = &cobra.Command{
	Use:   "reactions <message-url>",
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
		simple, _ := cmd.Flags().GetBool("simple")

		// リアクション一覧を取得
		reactions, err := client.GetReactions(messageURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: リアクションの取得に失敗しました: %v\n", err)
			os.Exit(1)
		}

		// スキントーンなどの修飾子を除去してリアクションを統合
		reactions = mergeReactions(reactions)

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

		// リアクション別にユーザー一覧を出力
		for _, reaction := range reactions {
			if simple {
				// シンプル形式：リアクション名のみ、改行区切りでユーザー一覧
				fmt.Fprintf(output, ":%s:\n", reaction.Name)
				for _, user := range reaction.Users {
					if email {
						fmt.Fprintln(output, user.Email)
					} else {
						fmt.Fprintln(output, user.Name)
					}
				}
				fmt.Fprintln(output) // 空行を追加
			} else {
				// 通常形式：リアクション名と人数、インデント付き
				fmt.Fprintf(output, ":%s: (%d人)\n", reaction.Name, len(reaction.Users))
				for _, user := range reaction.Users {
					if email {
						fmt.Fprintf(output, "  - %s\n", user.Email)
					} else {
						fmt.Fprintf(output, "  - %s\n", user.Name)
					}
				}
				fmt.Fprintln(output) // 空行を追加
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
	getCmd.AddCommand(getReactionsCmd)

	// reactions コマンドのフラグ（省略形用）
	reactionsCmd.Flags().StringP("filter", "f", "", "特定のリアクションのみをフィルタ（例: :参加します:）")
	reactionsCmd.Flags().BoolP("email", "e", false, "ユーザー名の代わりにメールアドレスを出力")
	reactionsCmd.Flags().StringP("output", "o", "", "結果をファイルに保存（例: reactions.txt）")
	reactionsCmd.Flags().BoolP("simple", "s", false, "シンプル形式で出力（改行のみで区切り、Googleカレンダーなどにコピーしやすい）")

	getReactionsCmd.Flags().StringP("filter", "f", "", "特定のリアクションのみをフィルタ（例: :参加します:）")
	getReactionsCmd.Flags().BoolP("email", "e", false, "ユーザー名の代わりにメールアドレスを出力")
	getReactionsCmd.Flags().StringP("output", "o", "", "結果をファイルに保存（例: reactions.txt）")
	getReactionsCmd.Flags().BoolP("simple", "s", false, "シンプル形式で出力（改行のみで区切り、Googleカレンダーなどにコピーしやすい）")
}

// mergeReactions merges reactions with skin tone modifiers into base reactions
func mergeReactions(reactions []slack.ReactionInfo) []slack.ReactionInfo {
	// 基本リアクション名をキーとして、ユーザーを統合
	reactionMap := make(map[string][]slack.UserInfo)

	for _, reaction := range reactions {
		baseName := normalizeReactionName(reaction.Name)

		// 既存のユーザーリストに追加
		if existingUsers, exists := reactionMap[baseName]; exists {
			reactionMap[baseName] = append(existingUsers, reaction.Users...)
		} else {
			reactionMap[baseName] = reaction.Users
		}
	}

	// マップからスライスに変換
	var merged []slack.ReactionInfo
	for baseName, users := range reactionMap {
		merged = append(merged, slack.ReactionInfo{
			Name:  baseName,
			Users: users,
		})
	}

	return merged
}

// normalizeReactionName removes skin tone and other modifiers from reaction names
func normalizeReactionName(name string) string {
	// スキントーンの修飾子を除去
	// 例: "+1::skin-tone-2:" -> "+1"
	// 例: "thumbsup::skin-tone-3:" -> "thumbsup"

	// ::skin-tone-X: パターンを除去
	if strings.Contains(name, "::skin-tone-") {
		parts := strings.Split(name, "::skin-tone-")
		if len(parts) > 0 {
			return parts[0]
		}
	}

	// その他の修飾子パターンも除去（必要に応じて追加）
	// 例: "::male:", "::female:", "::person:" など
	modifiers := []string{"::male:", "::female:", "::person:", "::man:", "::woman:"}
	for _, modifier := range modifiers {
		if strings.Contains(name, modifier) {
			parts := strings.Split(name, modifier)
			if len(parts) > 0 {
				return parts[0]
			}
		}
	}

	return name
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
