package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// saveToFile saves formatted content to a file
func saveToFile(content, filename, format string) error {
	// ファイル名が指定されていない場合はデフォルト名を生成
	if filename == "" {
		timestamp := time.Now().Format("20060102_150405")
		filename = fmt.Sprintf("slack_content_%s.md", timestamp)
	}

	// ファイル拡張子が指定されていない場合は.mdを追加
	if !strings.Contains(filename, ".") {
		filename += ".md"
	}

	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(filename)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("ディレクトリの作成に失敗しました: %v", err)
		}
	}

	// 形式に応じてコンテンツを変換
	var finalContent string
	switch strings.ToLower(format) {
	case "markdown", "md":
		finalContent = convertToMarkdown(content)
	default:
		finalContent = content
	}

	// ファイルに書き込み
	err := os.WriteFile(filename, []byte(finalContent), 0644)
	if err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました: %v", err)
	}

	return nil
}

// convertToMarkdown converts plain text content to markdown format
func convertToMarkdown(content string) string {
	// Markdownファイルでもプレーンテキストと同じ形式を維持
	// ヘッダーをMarkdownの見出しに変換し、取得日時も含める
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		// ヘッダー行をMarkdownの見出しに変換（取得日時も含める）
		if strings.Contains(line, "--- Slack") && strings.Contains(line, "の内容") {
			// 取得日時を抽出
			dateMatch := regexp.MustCompile(`\((\d{4}/\d{2}/\d{2} 取得)\)`).FindStringSubmatch(line)
			if len(dateMatch) > 1 {
				if strings.Contains(line, "スレッド") {
					result = append(result, fmt.Sprintf("# Slackスレッドの内容 (%s)", dateMatch[1]))
				} else if strings.Contains(line, "チャンネル") {
					result = append(result, fmt.Sprintf("# Slackチャンネルの内容 (%s)", dateMatch[1]))
				} else {
					result = append(result, fmt.Sprintf("# Slackの内容 (%s)", dateMatch[1]))
				}
			} else {
				if strings.Contains(line, "スレッド") {
					result = append(result, "# Slackスレッドの内容")
				} else if strings.Contains(line, "チャンネル") {
					result = append(result, "# Slackチャンネルの内容")
				} else {
					result = append(result, "# Slackの内容")
				}
			}
			continue
		}
		if strings.Contains(line, "--- ここまで ---") {
			// フッターは削除（重要ではないため）
			continue
		}

		// その他の行はそのまま追加（メッセージ形式は変更しない）
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}
