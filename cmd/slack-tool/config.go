package cmd

import (
	"fmt"
	"os"

	"github.com/shellme/slack-tool/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "設定の管理",
	Long:  "Slack APIトークンなどの設定を管理します。",
}

var configSetCmd = &cobra.Command{
	Use:   "set token <token>",
	Short: "Slack APIトークンを設定",
	Long:  "Slack APIのUser Tokenを設定ファイルに保存します。",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] != "token" {
			fmt.Fprintf(os.Stderr, "エラー: 無効な設定項目です。'token' を指定してください。\n")
			os.Exit(1)
		}

		token := args[1]

		// トークンの形式を検証
		if err := config.ValidateToken(token); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}

		// 設定マネージャーを作成
		cm := config.NewConfigManager()

		// 現在の設定を読み込み
		cfg, err := cm.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: 設定の読み込みに失敗しました: %v\n", err)
			os.Exit(1)
		}

		// トークンを設定
		cfg.SlackToken = token

		// 設定を保存
		if err := cm.Save(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: 設定の保存に失敗しました: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Slack APIトークンが正常に設定されました。")
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "現在の設定を表示",
	Long:  "現在設定されているSlack APIトークンをマスクして表示します。",
	Run: func(cmd *cobra.Command, args []string) {
		// 設定マネージャーを作成
		cm := config.NewConfigManager()

		// 設定を読み込み
		cfg, err := cm.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: 設定の読み込みに失敗しました: %v\n", err)
			os.Exit(1)
		}

		// 設定ファイルのパスを表示
		fmt.Printf("設定ファイル: %s\n", cm.GetConfigPath())

		// トークンの状態を表示
		if cfg.SlackToken == "" {
			fmt.Println("Slack APIトークン: 未設定")
			fmt.Println("\nトークンを設定するには:")
			fmt.Println("  slack-tool config set token \"xoxp-xxxxxxxxxxxxxx-xxxxxxxx\"")
		} else {
			maskedToken := config.MaskToken(cfg.SlackToken)
			fmt.Printf("Slack APIトークン: %s\n", maskedToken)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
}
