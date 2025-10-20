package slack

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

// Formatter handles formatting of Slack messages for output
type Formatter struct {
	client     *Client
	users      map[string]*slack.User      // ユーザー情報のキャッシュ
	usergroups map[string]*slack.UserGroup // サブチーム情報のキャッシュ
}

// NewFormatter creates a new formatter
func NewFormatter(client *Client) *Formatter {
	return &Formatter{
		client:     client,
		users:      make(map[string]*slack.User),
		usergroups: make(map[string]*slack.UserGroup),
	}
}

// FormatThread formats a thread of messages for output
func (f *Formatter) FormatThread(messages []slack.Message) (string, error) {
	if len(messages) == 0 {
		return "", fmt.Errorf("メッセージがありません")
	}

	// JSTタイムゾーンを設定
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "", fmt.Errorf("タイムゾーンの設定に失敗しました: %v", err)
	}

	var result strings.Builder

	// ヘッダーを追加
	result.WriteString("--- Slackスレッドの内容 (")
	result.WriteString(time.Now().In(jst).Format("2006/01/02 取得"))
	result.WriteString(") ---\n\n")

	// 各メッセージをフォーマット
	for _, msg := range messages {
		formatted, err := f.formatMessage(msg, jst)
		if err != nil {
			return "", fmt.Errorf("メッセージのフォーマットに失敗しました: %v", err)
		}

		result.WriteString(formatted)
		result.WriteString("\n\n") // メッセージ間に空行を追加
	}

	// フッターを追加
	result.WriteString("--- ここまで ---")

	return result.String(), nil
}

// FormatChannel formats channel messages for output
func (f *Formatter) FormatChannel(messages []slack.Message, channelName string) (string, error) {
	if len(messages) == 0 {
		return "", fmt.Errorf("メッセージがありません")
	}

	// JSTタイムゾーンを設定
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "", fmt.Errorf("タイムゾーンの設定に失敗しました: %v", err)
	}

	var result strings.Builder

	// ヘッダーを追加
	result.WriteString("--- Slackチャンネルの内容 (")
	result.WriteString(time.Now().In(jst).Format("2006/01/02 取得"))
	result.WriteString(") ---\n")
	if channelName != "" {
		result.WriteString(fmt.Sprintf("チャンネル: #%s\n", channelName))
	}
	result.WriteString("\n")

	// メッセージをスレッド構造でフォーマット
	formatted, err := f.formatChannelWithThreads(messages, jst)
	if err != nil {
		return "", fmt.Errorf("メッセージのフォーマットに失敗しました: %v", err)
	}

	result.WriteString(formatted)

	// フッターを追加
	result.WriteString("--- ここまで ---")

	return result.String(), nil
}

// formatChannelWithThreads formats channel messages with thread structure
func (f *Formatter) formatChannelWithThreads(messages []slack.Message, jst *time.Location) (string, error) {
	// メインメッセージとスレッド返信を分離
	mainMessages := make(map[string]slack.Message)
	threadReplies := make(map[string][]slack.Message)

	for _, msg := range messages {
		if msg.ThreadTimestamp == "" || msg.ThreadTimestamp == msg.Timestamp {
			// メインメッセージ
			mainMessages[msg.Timestamp] = msg
		} else {
			// スレッド返信
			threadReplies[msg.ThreadTimestamp] = append(threadReplies[msg.ThreadTimestamp], msg)
		}
	}

	var result strings.Builder

	// メインメッセージを時系列で処理
	for _, msg := range messages {
		if msg.ThreadTimestamp == "" || msg.ThreadTimestamp == msg.Timestamp {
			// メインメッセージをフォーマット
			formatted, err := f.formatMessage(msg, jst)
			if err != nil {
				return "", fmt.Errorf("メッセージのフォーマットに失敗しました: %v", err)
			}
			result.WriteString(formatted)
			result.WriteString("\n")

			// このメッセージのスレッド返信があるかチェック
			if replies, exists := threadReplies[msg.Timestamp]; exists {
				// スレッド返信を時系列でソート
				sortedReplies := f.sortMessagesByTimestamp(replies)

				for _, reply := range sortedReplies {
					replyFormatted, err := f.formatMessage(reply, jst)
					if err != nil {
						return "", fmt.Errorf("スレッド返信のフォーマットに失敗しました: %v", err)
					}

					// スレッド返信としてインデントして表示
					lines := strings.Split(replyFormatted, "\n")
					for i, line := range lines {
						if i == 0 {
							result.WriteString(fmt.Sprintf("  └─ %s\n", line))
						} else {
							result.WriteString(fmt.Sprintf("     %s\n", line))
						}
					}
				}
			}

			result.WriteString("\n") // メッセージ間に空行を追加
		}
	}

	return result.String(), nil
}

// sortMessagesByTimestamp sorts messages by timestamp
func (f *Formatter) sortMessagesByTimestamp(messages []slack.Message) []slack.Message {
	// バブルソートでタイムスタンプ順にソート
	for i := 0; i < len(messages)-1; i++ {
		for j := 0; j < len(messages)-1-i; j++ {
			if messages[j].Timestamp > messages[j+1].Timestamp {
				messages[j], messages[j+1] = messages[j+1], messages[j]
			}
		}
	}
	return messages
}

// formatMessage formats a single message
func (f *Formatter) formatMessage(msg slack.Message, jst *time.Location) (string, error) {
	// ユーザー情報を取得（キャッシュから、またはAPIから）
	user, err := f.getUserInfo(msg.User)
	var username string
	if err != nil {
		// ユーザー情報が取得できない場合はユーザーIDをそのまま使用
		username = "@" + msg.User
	} else {
		// ユーザー名を取得（@以降の名前）
		username = f.getUsername(user)
	}

	// タイムスタンプをJSTに変換
	timestamp, err := f.parseTimestamp(msg.Timestamp)
	if err != nil {
		return "", fmt.Errorf("タイムスタンプの解析に失敗しました: %v", err)
	}

	jstTime := timestamp.In(jst)
	timeStr := jstTime.Format("2006-01-02 15:04:05")

	// メッセージテキストをクリーンアップ
	text := f.cleanMessageText(msg.Text)

	// フォーマット: [YYYY-MM-DD HH:MM:SS][@username]: 本文
	return fmt.Sprintf("[%s][%s]:\n%s", timeStr, username, text), nil
}

// getUserInfo gets user information, using cache if available
func (f *Formatter) getUserInfo(userID string) (*slack.User, error) {
	// キャッシュをチェック
	if user, exists := f.users[userID]; exists {
		return user, nil
	}

	// APIから取得
	user, err := f.client.GetUserInfo(userID)
	if err != nil {
		return nil, err
	}

	// キャッシュに保存
	f.users[userID] = user

	return user, nil
}

// getUserGroupInfo gets user group (subteam) information, using cache if available
func (f *Formatter) getUserGroupInfo(groupID string) (*slack.UserGroup, error) {
	// キャッシュをチェック
	if group, exists := f.usergroups[groupID]; exists {
		return group, nil
	}

	// サブチーム情報がキャッシュにない場合は、全サブチーム情報を取得してキャッシュに保存
	if len(f.usergroups) == 0 {
		usergroups, err := f.client.GetUserGroups()
		if err != nil {
			return nil, err
		}

		// 全サブチーム情報をキャッシュに保存
		for _, group := range usergroups {
			f.usergroups[group.ID] = &group
		}
	}

	// キャッシュから取得
	if group, exists := f.usergroups[groupID]; exists {
		return group, nil
	}

	// 見つからない場合はnilを返す（エラーではない）
	return nil, nil
}

// getUsername extracts the @username from user information
func (f *Formatter) getUsername(user *slack.User) string {
	// ユーザー名の優先順位: Name > RealName > ID
	if user.Name != "" {
		return "@" + user.Name
	}
	if user.RealName != "" {
		return "@" + user.RealName
	}
	return "@" + user.ID
}

// parseTimestamp parses Slack timestamp string to time.Time
func (f *Formatter) parseTimestamp(timestamp string) (time.Time, error) {
	// Slackのタイムスタンプは "1234567890.123456" 形式
	parts := strings.Split(timestamp, ".")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("無効なタイムスタンプ形式: %s", timestamp)
	}

	// 秒部分を解析
	sec, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("タイムスタンプの秒部分の解析に失敗: %v", err)
	}

	// マイクロ秒部分を解析（最大6桁）
	microsecStr := parts[1]
	if len(microsecStr) > 6 {
		microsecStr = microsecStr[:6]
	}

	// 6桁にパディング
	for len(microsecStr) < 6 {
		microsecStr += "0"
	}

	microsec, err := strconv.ParseInt(microsecStr, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("タイムスタンプのマイクロ秒部分の解析に失敗: %v", err)
	}

	// Unix時間からtime.Timeに変換
	return time.Unix(sec, microsec*1000), nil
}

// cleanMessageText cleans up message text for output
func (f *Formatter) cleanMessageText(text string) string {
	// 基本的なクリーンアップ
	cleaned := strings.TrimSpace(text)

	// 空のメッセージの場合は適切なメッセージを表示
	if cleaned == "" {
		return "(メッセージの内容がありません)"
	}

	// 改行を正規化
	cleaned = strings.ReplaceAll(cleaned, "\r\n", "\n")
	cleaned = strings.ReplaceAll(cleaned, "\r", "\n")

	// メンションを変換
	cleaned = f.convertMentions(cleaned)

	return cleaned
}

// convertMentions converts Slack mention IDs to readable usernames
func (f *Formatter) convertMentions(text string) string {
	// ユーザーメンションの正規表現: <@U08KTGLLCLU>
	userMentionRegex := regexp.MustCompile(`<@([UW][A-Z0-9]+)>`)

	// サブチームメンションの正規表現: <!subteam^S025PC88BJ5>
	subteamMentionRegex := regexp.MustCompile(`<!subteam\^([A-Z0-9]+)>`)

	// チャンネルメンションの正規表現: <#C12345678|channel-name>
	channelMentionRegex := regexp.MustCompile(`<#([A-Z0-9]+)(\|[^>]+)?>`)

	// ユーザーメンションを変換
	text = userMentionRegex.ReplaceAllStringFunc(text, func(match string) string {
		// <@U08KTGLLCLU> から U08KTGLLCLU を抽出
		userID := userMentionRegex.FindStringSubmatch(match)[1]

		// ユーザー情報を取得
		user, err := f.getUserInfo(userID)
		if err != nil {
			// ユーザー情報が取得できない場合は元のIDを表示
			return "@" + userID
		}

		// ユーザー名を取得（@は除く）
		username := f.getUsername(user)
		username = strings.TrimPrefix(username, "@")

		return "@" + username
	})

	// サブチームメンションを変換
	text = subteamMentionRegex.ReplaceAllStringFunc(text, func(match string) string {
		// <!subteam^S025PC88BJ5> から S025PC88BJ5 を抽出
		subteamID := subteamMentionRegex.FindStringSubmatch(match)[1]

		// サブチーム情報を取得
		group, err := f.getUserGroupInfo(subteamID)
		if err != nil {
			// エラーの場合はIDをそのまま表示
			return "@" + subteamID
		}

		if group != nil && group.Handle != "" {
			// ハンドル名が取得できた場合
			return "@" + group.Handle
		}

		// ハンドル名が取得できない場合はIDをそのまま表示
		return "@" + subteamID
	})

	// チャンネルメンションを変換
	text = channelMentionRegex.ReplaceAllStringFunc(text, func(match string) string {
		// <#C12345678|channel-name> から C12345678 と channel-name を抽出
		matches := channelMentionRegex.FindStringSubmatch(match)
		channelID := matches[1]
		channelName := ""

		// チャンネル名が含まれている場合
		if len(matches) > 2 && matches[2] != "" {
			channelName = matches[2][1:] // |を除去
		}

		if channelName != "" {
			return "#" + channelName
		}

		return "#" + channelID
	})

	return text
}
