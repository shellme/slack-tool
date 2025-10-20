package cmd

import (
	"github.com/spf13/cobra"
)

// getThreadShortcutCmd は `get` コマンドの省略形で、`get thread` と同じ動作をします
var getThreadShortcutCmd = &cobra.Command{
	Use:   "get-thread <slack-thread-url>",
	Short: "スレッドの内容を取得・整形（get thread の省略形）",
	Aliases: []string{"get"},
	Long: `指定されたSlackスレッドのURLから会話内容を取得し、
AIへの入力に適した人間が読みやすいプレーンテキスト形式で整形して表示します。

これは 'get thread' コマンドの省略形です。

例:
  slack-tool get "https://your-workspace.slack.com/archives/C12345678/p1234567890123456"
  slack-tool get "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --output thread.md
  slack-tool get "https://your-workspace.slack.com/archives/C12345678/p1234567890123456" --output thread.md --format markdown`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// get thread コマンドと同じ処理を実行
		getThreadCmd.Run(cmd, args)
	},
}

// postMessageShortcutCmd は `post` コマンドの省略形で、`post message` と同じ動作をします
var postMessageShortcutCmd = &cobra.Command{
	Use:   "post-message <message>",
	Short: "メッセージを投稿（post message の省略形）",
	Aliases: []string{"post"},
	Long: `指定されたSlackチャンネルにメッセージを投稿します。

これは 'post message' コマンドの省略形です。

例:
  slack-tool post "Hello, world!" --channel C12345678
  slack-tool post "This is a test message" --channel C12345678 --thread 1234567890.123456
  slack-tool post "Hey @john, can you review this?" --channel C12345678
  slack-tool post "スレッド返信です" --thread-url "https://workspace.slack.com/archives/C12345678/p1234567890123456"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// post message コマンドと同じ処理を実行
		postMessageCmd.Run(cmd, args)
	},
}

// getReactionsShortcutCmd は `reactions` コマンドの省略形で、`get reactions` と同じ動作をします
var getReactionsShortcutCmd = &cobra.Command{
	Use:   "get-reactions <message-url>",
	Short: "指定した投稿のリアクション一覧を取得（get reactions の省略形）",
	Aliases: []string{"reactions"},
	Long: `指定したSlack投稿のリアクション一覧を取得します。

これは 'get reactions' コマンドの省略形です。

例:
  slack-tool reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456"
  slack-tool reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --filter ":参加します:"
  slack-tool reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --email
  slack-tool reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --output reactions.txt
  slack-tool reactions "https://workspace.slack.com/archives/C12345678/p1234567890123456" --filter "承知_しました" --email --output participants.txt`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// get reactions コマンドと同じ処理を実行
		getReactionsCmd.Run(cmd, args)
	},
}

func init() {
	// 省略コマンドをルートコマンドに追加
	rootCmd.AddCommand(getThreadShortcutCmd)
	rootCmd.AddCommand(postMessageShortcutCmd)
	rootCmd.AddCommand(getReactionsShortcutCmd)

	// 各省略コマンドに元のコマンドと同じフラグを追加
	getThreadShortcutCmd.Flags().StringP("output", "o", "", "出力ファイル名を指定（例: thread.md, thread.txt）。拡張子で形式を自動判定")
	getThreadShortcutCmd.Flags().StringP("format", "f", "text", "出力形式を指定（text / markdown）。指定があれば拡張子より優先")

	postMessageShortcutCmd.Flags().StringP("channel", "c", "", "投稿先のチャンネルIDまたはURL")
	postMessageShortcutCmd.Flags().StringP("thread", "t", "", "スレッド返信する場合のタイムスタンプ")
	postMessageShortcutCmd.Flags().StringP("thread-url", "u", "", "スレッド返信する場合のスレッドURL")

	getReactionsShortcutCmd.Flags().StringP("filter", "f", "", "特定のリアクションのみをフィルタ（例: :参加します:）")
	getReactionsShortcutCmd.Flags().BoolP("email", "e", false, "ユーザー名の代わりにメールアドレスを出力")
	getReactionsShortcutCmd.Flags().StringP("output", "o", "", "結果をファイルに保存（例: reactions.txt）")
}
