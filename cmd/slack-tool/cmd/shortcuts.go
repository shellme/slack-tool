package cmd

import (
	"github.com/spf13/cobra"
)

// getShortcutCmd は `get` コマンドの省略形で、`get thread` と同じ動作をします
var getShortcutCmd = &cobra.Command{
	Use:   "get <slack-thread-url>",
	Short: "スレッドの内容を取得・整形（get thread の省略形）",
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
		getCmd.Run(cmd, args)
	},
}

// postShortcutCmd は `post` コマンドの省略形で、`post message` と同じ動作をします
var postShortcutCmd = &cobra.Command{
	Use:   "post <message>",
	Short: "メッセージを投稿（post message の省略形）",
	Long: `指定されたSlackチャンネルにメッセージを投稿します。

これは 'post message' コマンドの省略形です。

例:
  slack-tool post "Hello, world!" --channel C12345678
  slack-tool post "This is a test message" --channel C12345678 --thread 1234567890.123456
  slack-tool post "Hey @john, can you review this?" --channel C12345678
  slack-tool post "スレッド返信です" --thread-url "https://workspace.slack.com/archives/C12345678/p1234567890123456"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// post コマンドと同じ処理を実行
		postCmd.Run(cmd, args)
	},
}

// reactionsShortcutCmd は `reactions` コマンドの省略形で、`get reactions` と同じ動作をします
var reactionsShortcutCmd = &cobra.Command{
	Use:   "reactions <message-url>",
	Short: "指定した投稿のリアクション一覧を取得（get reactions の省略形）",
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
		reactionsCmd.Run(cmd, args)
	},
}

func init() {
	// 省略コマンドをルートコマンドに追加
	rootCmd.AddCommand(getShortcutCmd)
	rootCmd.AddCommand(postShortcutCmd)
	rootCmd.AddCommand(reactionsShortcutCmd)

	// 各省略コマンドに元のコマンドと同じフラグを追加
	getShortcutCmd.Flags().StringP("output", "o", "", "出力ファイル名を指定（例: thread.md, thread.txt）。拡張子で形式を自動判定")
	getShortcutCmd.Flags().StringP("format", "f", "text", "出力形式を指定（text / markdown）。指定があれば拡張子より優先")

	postShortcutCmd.Flags().StringP("channel", "c", "", "投稿先のチャンネルIDまたはURL")
	postShortcutCmd.Flags().StringP("thread", "t", "", "スレッド返信する場合のタイムスタンプ")
	postShortcutCmd.Flags().StringP("thread-url", "u", "", "スレッド返信する場合のスレッドURL")

	reactionsShortcutCmd.Flags().StringP("filter", "f", "", "特定のリアクションのみをフィルタ（例: :参加します:）")
	reactionsShortcutCmd.Flags().BoolP("email", "e", false, "ユーザー名の代わりにメールアドレスを出力")
	reactionsShortcutCmd.Flags().StringP("output", "o", "", "結果をファイルに保存（例: reactions.txt）")
}
