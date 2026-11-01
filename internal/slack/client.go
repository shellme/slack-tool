package slack

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

// historyPageSize is the page size for conversations.history (Slack recommends 200 or fewer)
const historyPageSize = 200

// maxRateLimitRetries is the number of retries after a rate limit error
const maxRateLimitRetries = 5

// retry calls fn, waiting and retrying when Slack returns a rate limit error
func retry[T any](fn func() (T, error)) (T, error) {
	for attempt := 0; ; attempt++ {
		result, err := fn()

		var rateLimited *slack.RateLimitedError
		if !errors.As(err, &rateLimited) || attempt >= maxRateLimitRetries {
			return result, err
		}

		fmt.Fprintf(os.Stderr, "情報: APIレート制限のため %s 待機します\n", rateLimited.RetryAfter)
		time.Sleep(rateLimited.RetryAfter)
	}
}

// Client wraps the Slack API client
type Client struct {
	api *slack.Client
}

// NewClient creates a new Slack client
func NewClient(token string) *Client {
	api := slack.New(token)
	return &Client{api: api}
}

// GetThreadReplies fetches all replies in a thread
func (c *Client) GetThreadReplies(channelID, timestamp string) ([]slack.Message, error) {
	var messages []slack.Message
	cursor := ""

	for {
		// conversations.replies APIを呼び出し
		params := &slack.GetConversationRepliesParameters{
			ChannelID: channelID,
			Timestamp: timestamp,
			Cursor:    cursor,
			Inclusive: true, // 指定されたタイムスタンプのメッセージも含める
		}

		var hasMore bool
		var nextCursor string
		page, err := retry(func() ([]slack.Message, error) {
			var msgs []slack.Message
			var err error
			msgs, hasMore, nextCursor, err = c.api.GetConversationReplies(params)
			return msgs, err
		})
		if err != nil {
			return nil, c.handleAPIError(err)
		}

		messages = append(messages, page...)
		if !hasMore || nextCursor == "" {
			return messages, nil
		}
		cursor = nextCursor
	}
}

// GetUserInfo fetches user information by user ID
func (c *Client) GetUserInfo(userID string) (*slack.User, error) {
	user, err := retry(func() (*slack.User, error) {
		return c.api.GetUserInfo(userID)
	})
	if err != nil {
		return nil, c.handleAPIError(err)
	}

	return user, nil
}

// GetUserGroups fetches all user groups (subteams) information
func (c *Client) GetUserGroups() ([]slack.UserGroup, error) {
	usergroups, err := retry(func() ([]slack.UserGroup, error) {
		return c.api.GetUserGroups()
	})
	if err != nil {
		return nil, c.handleAPIError(err)
	}

	return usergroups, nil
}

// GetChannelHistory fetches channel history
func (c *Client) GetChannelHistory(channelID string, limit int) ([]slack.Message, error) {
	return c.fetchHistory(channelID, limit, "", "")
}

// GetChannelHistoryWithThreads fetches channel history including thread replies
func (c *Client) GetChannelHistoryWithThreads(channelID string, limit int) ([]slack.Message, error) {
	return c.GetChannelHistoryWithThreadsInRange(channelID, limit, "", "")
}

// GetChannelHistoryWithThreadsInRange fetches channel history including thread replies with date range
func (c *Client) GetChannelHistoryWithThreadsInRange(channelID string, limit int, oldest, latest string) ([]slack.Message, error) {
	// 期間指定がある場合は設定
	var oldestTS, latestTS string
	var err error
	if oldest != "" {
		oldestTS, err = c.parseTimestamp(oldest)
		if err != nil {
			return nil, fmt.Errorf("oldest の日時解析に失敗しました: %v", err)
		}
	}
	if latest != "" {
		latestTS, err = c.parseTimestamp(latest)
		if err != nil {
			return nil, fmt.Errorf("latest の日時解析に失敗しました: %v", err)
		}
	}

	messages, err := c.fetchHistory(channelID, limit, oldestTS, latestTS)
	if err != nil {
		return nil, err
	}

	// スレッドの返信も取得
	// チャンネルにも投稿された返信は履歴と重複するため除外する
	var allMessages []slack.Message
	seen := make(map[string]bool)
	add := func(msg slack.Message) {
		if !seen[msg.Timestamp] {
			seen[msg.Timestamp] = true
			allMessages = append(allMessages, msg)
		}
	}

	for _, msg := range messages {
		add(msg)

		// 親メッセージの場合のみスレッドの返信を取得
		if msg.ThreadTimestamp != "" && msg.ThreadTimestamp == msg.Timestamp {
			threadReplies, err := c.GetThreadReplies(channelID, msg.ThreadTimestamp)
			if err != nil {
				// エラーが発生してもメインメッセージは含める
				fmt.Fprintf(os.Stderr, "警告: スレッド %s の返信を取得できませんでした: %v\n", msg.ThreadTimestamp, err)
				continue
			}

			for _, reply := range threadReplies {
				add(reply)
			}
		}
	}

	return allMessages, nil
}

// fetchHistory fetches up to limit messages, following pagination
func (c *Client) fetchHistory(channelID string, limit int, oldest, latest string) ([]slack.Message, error) {
	var messages []slack.Message
	cursor := ""

	for {
		pageSize := historyPageSize
		if limit-len(messages) < pageSize {
			pageSize = limit - len(messages)
		}

		params := &slack.GetConversationHistoryParameters{
			ChannelID: channelID,
			Cursor:    cursor,
			Limit:     pageSize,
			Oldest:    oldest,
			Latest:    latest,
		}

		response, err := retry(func() (*slack.GetConversationHistoryResponse, error) {
			return c.api.GetConversationHistory(params)
		})
		if err != nil {
			return nil, c.handleAPIError(err)
		}

		messages = append(messages, response.Messages...)
		cursor = response.ResponseMetaData.NextCursor
		if !response.HasMore || cursor == "" || len(messages) >= limit {
			return messages, nil
		}
	}
}

// GetChannelInfo fetches channel information
func (c *Client) GetChannelInfo(channelID string) (*slack.Channel, error) {
	channel, err := retry(func() (*slack.Channel, error) {
		return c.api.GetConversationInfo(&slack.GetConversationInfoInput{
			ChannelID: channelID,
		})
	})
	if err != nil {
		return nil, c.handleAPIError(err)
	}

	return channel, nil
}

// PostMessage posts a message to a Slack channel
func (c *Client) PostMessage(channelID, text string) error {
	return c.PostMessageWithOptions(channelID, text)
}

// PostMessageWithOptions posts a message with additional options
func (c *Client) PostMessageWithOptions(channelID, text string, options ...slack.MsgOption) error {
	msgOptions := []slack.MsgOption{
		slack.MsgOptionText(text, false),
	}
	msgOptions = append(msgOptions, options...)

	// レート制限時は投稿されないため再試行しても重複しない
	_, err := retry(func() (string, error) {
		_, ts, err := c.api.PostMessage(channelID, msgOptions...)
		return ts, err
	})
	if err != nil {
		return c.handleAPIError(err)
	}
	return nil
}

// PostThreadReply posts a reply to a thread
func (c *Client) PostThreadReply(channelID, text, threadTimestamp string) error {
	// Slack API ドキュメントに基づく正しい実装
	return c.PostMessageWithOptions(channelID, text,
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			ThreadTimestamp: threadTimestamp,
		}))
}

// PostThreadReplyByURL posts a reply to a thread using a thread URL
func (c *Client) PostThreadReplyByURL(text, threadURL string) error {
	// スレッドURLを解析
	threadInfo, err := ParseThreadURL(threadURL)
	if err != nil {
		return fmt.Errorf("スレッドURLの解析に失敗しました: %v", err)
	}

	// メッセージ情報を取得して正確なタイムスタンプを確認
	msg, err := c.GetMessageInfo(threadInfo.ChannelID, threadInfo.Timestamp, threadInfo.ThreadTimestamp)
	if err != nil {
		return fmt.Errorf("メッセージ情報の取得に失敗しました: %v", err)
	}

	// 返信のURLが指定された場合は親スレッドに返信する
	threadTimestamp := msg.Timestamp
	if msg.ThreadTimestamp != "" {
		threadTimestamp = msg.ThreadTimestamp
	}

	return c.PostThreadReply(threadInfo.ChannelID, text, threadTimestamp)
}

// GetMessageInfo gets a message by timestamp. threadTimestamp is required to find a thread reply
func (c *Client) GetMessageInfo(channelID, timestamp, threadTimestamp string) (*slack.Message, error) {
	// スレッド返信でなければタイムスタンプを指定して直接取得
	if threadTimestamp == "" || threadTimestamp == timestamp {
		params := &slack.GetConversationHistoryParameters{
			ChannelID: channelID,
			Oldest:    timestamp,
			Latest:    timestamp,
			Inclusive: true,
			Limit:     1,
		}

		response, err := retry(func() (*slack.GetConversationHistoryResponse, error) {
			return c.api.GetConversationHistory(params)
		})
		if err != nil {
			return nil, c.handleAPIError(err)
		}

		for _, msg := range response.Messages {
			if msg.Timestamp == timestamp {
				return &msg, nil
			}
		}

		return nil, fmt.Errorf("指定されたタイムスタンプのメッセージが見つかりませんでした: %s（スレッド返信の場合は thread_ts 付きのURLを指定してください）", timestamp)
	}

	// スレッド内のメッセージを検索
	replies, err := c.GetThreadReplies(channelID, threadTimestamp)
	if err != nil {
		return nil, err
	}

	for _, reply := range replies {
		if reply.Timestamp == timestamp {
			return &reply, nil
		}
	}

	return nil, fmt.Errorf("指定されたタイムスタンプのメッセージが見つかりませんでした: %s", timestamp)
}

// ReactionInfo contains information about a reaction
type ReactionInfo struct {
	Name  string     `json:"name"`
	Users []UserInfo `json:"users"`
}

// UserInfo contains basic user information
type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetReactions gets reactions for a specific message
func (c *Client) GetReactions(messageURL string) ([]ReactionInfo, error) {
	// メッセージURLを解析
	threadInfo, err := ParseThreadURL(messageURL)
	if err != nil {
		return nil, fmt.Errorf("メッセージURLの解析に失敗しました: %v", err)
	}

	// メッセージ情報を取得
	_, err = c.GetMessageInfo(threadInfo.ChannelID, threadInfo.Timestamp, threadInfo.ThreadTimestamp)
	if err != nil {
		return nil, fmt.Errorf("メッセージ情報の取得に失敗しました: %v", err)
	}

	// リアクション情報を取得
	itemRef := slack.ItemRef{
		Channel:   threadInfo.ChannelID,
		Timestamp: threadInfo.Timestamp,
	}

	item, err := retry(func() (slack.ReactedItem, error) {
		return c.api.GetReactions(itemRef, slack.GetReactionsParameters{})
	})
	if err != nil {
		return nil, c.handleAPIError(err)
	}

	// リアクション情報を変換
	var reactionInfos []ReactionInfo
	for _, reaction := range item.Reactions {
		// ユーザー情報を取得
		var users []UserInfo
		for _, userID := range reaction.Users {
			user, err := c.GetUserInfo(userID)
			if err != nil {
				// ユーザー情報が取得できない場合もリアクションした事実は残す
				fmt.Fprintf(os.Stderr, "警告: ユーザー %s の情報を取得できませんでした: %v\n", userID, err)
				users = append(users, UserInfo{ID: userID})
				continue
			}
			users = append(users, UserInfo{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Profile.Email,
			})
		}

		reactionInfos = append(reactionInfos, ReactionInfo{
			Name:  reaction.Name,
			Users: users,
		})
	}

	return reactionInfos, nil
}

// handleAPIError converts Slack API errors to user-friendly messages
func (c *Client) handleAPIError(err error) error {
	if err == nil {
		return nil
	}

	var rateLimited *slack.RateLimitedError
	if errors.As(err, &rateLimited) {
		return fmt.Errorf("APIレート制限に達しました。しばらく待ってから再試行してください")
	}

	// エラーメッセージをチェックして適切な日本語メッセージに変換
	errStr := err.Error()

	switch {
	case contains(errStr, "invalid_auth"):
		return fmt.Errorf("認証に失敗しました。トークンが無効または期限切れです")
	case contains(errStr, "account_inactive"):
		return fmt.Errorf("アカウントが無効です")
	case contains(errStr, "token_revoked"):
		return fmt.Errorf("トークンが取り消されました")
	case contains(errStr, "not_authed"):
		return fmt.Errorf("認証されていません")
	case contains(errStr, "channel_not_found"):
		return fmt.Errorf("チャンネルが見つかりません")
	case contains(errStr, "thread_not_found"):
		return fmt.Errorf("スレッドが見つかりません")
	case contains(errStr, "not_in_channel"):
		return fmt.Errorf("このチャンネルにアクセスする権限がありません")
	case contains(errStr, "rate_limited"):
		return fmt.Errorf("APIレート制限に達しました。しばらく待ってから再試行してください")
	case contains(errStr, "timeout"):
		return fmt.Errorf("ネットワークタイムアウトが発生しました")
	case contains(errStr, "no_network"):
		return fmt.Errorf("ネットワーク接続エラーが発生しました")
	default:
		return fmt.Errorf("Slack APIエラー: %v", err)
	}
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr))))
}

// containsSubstring performs a simple substring search
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// parseTimestamp converts various date formats to Slack timestamp format
func (c *Client) parseTimestamp(dateStr string) (string, error) {
	// 既にUnixタイムスタンプ形式（数字のみ）の場合はそのまま返す
	if timestamp, err := strconv.ParseFloat(dateStr, 64); err == nil {
		// Slack APIは秒単位のタイムスタンプを期待する
		return fmt.Sprintf("%.0f", timestamp), nil
	}

	// 日時文字列を解析するためのフォーマット一覧
	formats := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02T15:04:05",       // ローカル時刻
		"2006-01-02 15:04:05",       // スペース区切り
		"2006-01-02",                // 日付のみ
	}

	for _, format := range formats {
		// タイムゾーン指定がない場合はローカル時刻として扱う
		if t, err := time.ParseInLocation(format, dateStr, time.Local); err == nil {
			// Unixタイムスタンプに変換（秒単位）
			return fmt.Sprintf("%.0f", float64(t.Unix())), nil
		}
	}

	return "", fmt.Errorf("サポートされていない日時形式です: %s", dateStr)
}

// TestConnection tests the Slack API connection
func (c *Client) TestConnection() error {
	// auth.test APIを呼び出して接続をテスト
	_, err := c.api.AuthTest()
	if err != nil {
		return c.handleAPIError(err)
	}

	return nil
}
