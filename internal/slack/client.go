package slack

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

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
	// conversations.replies APIを呼び出し
	params := &slack.GetConversationRepliesParameters{
		ChannelID: channelID,
		Timestamp: timestamp,
		Inclusive: true, // 指定されたタイムスタンプのメッセージも含める
	}

	messages, _, _, err := c.api.GetConversationReplies(params)
	if err != nil {
		return nil, c.handleAPIError(err)
	}

	return messages, nil
}

// GetUserInfo fetches user information by user ID
func (c *Client) GetUserInfo(userID string) (*slack.User, error) {
	user, err := c.api.GetUserInfo(userID)
	if err != nil {
		return nil, c.handleAPIError(err)
	}

	return user, nil
}

// GetUserGroups fetches all user groups (subteams) information
func (c *Client) GetUserGroups() ([]slack.UserGroup, error) {
	usergroups, err := c.api.GetUserGroups()
	if err != nil {
		return nil, c.handleAPIError(err)
	}

	return usergroups, nil
}

// GetChannelHistory fetches channel history
func (c *Client) GetChannelHistory(channelID string, limit int) ([]slack.Message, error) {
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     limit,
	}

	messages, err := c.api.GetConversationHistory(params)
	if err != nil {
		return nil, c.handleAPIError(err)
	}

	return messages.Messages, nil
}

// GetChannelHistoryWithThreads fetches channel history including thread replies
func (c *Client) GetChannelHistoryWithThreads(channelID string, limit int) ([]slack.Message, error) {
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     limit,
	}

	messages, err := c.api.GetConversationHistory(params)
	if err != nil {
		return nil, c.handleAPIError(err)
	}

	// スレッドの返信も取得
	var allMessages []slack.Message
	threadCount := 0

	for _, msg := range messages.Messages {
		allMessages = append(allMessages, msg)

		// スレッドの返信があるかチェック
		// ThreadTimestampが空でない場合、そのメッセージはスレッドに関連している
		if msg.ThreadTimestamp != "" {
			threadCount++

			// スレッドの返信を取得（ThreadTimestampを使用）
			threadReplies, err := c.GetThreadReplies(channelID, msg.ThreadTimestamp)
			if err != nil {
				// エラーが発生してもメインメッセージは含める
				continue
			}

			// スレッドの返信を追加（メインメッセージは除外）
			for _, reply := range threadReplies {
				if reply.Timestamp != msg.Timestamp {
					allMessages = append(allMessages, reply)
				}
			}
		}
	}

	return allMessages, nil
}

// GetChannelInfo fetches channel information
func (c *Client) GetChannelInfo(channelID string) (*slack.Channel, error) {
	channel, err := c.api.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID: channelID,
	})
	if err != nil {
		return nil, c.handleAPIError(err)
	}

	return channel, nil
}

// PostMessage posts a message to a Slack channel
func (c *Client) PostMessage(channelID, text string) error {
	_, _, err := c.api.PostMessage(channelID, slack.MsgOptionText(text, false))
	if err != nil {
		return c.handleAPIError(err)
	}
	return nil
}

// PostMessageWithOptions posts a message with additional options
func (c *Client) PostMessageWithOptions(channelID, text string, options ...slack.MsgOption) error {
	msgOptions := []slack.MsgOption{
		slack.MsgOptionText(text, false),
	}
	msgOptions = append(msgOptions, options...)

	_, _, err := c.api.PostMessage(channelID, msgOptions...)
	if err != nil {
		return c.handleAPIError(err)
	}
	return nil
}

// PostThreadReply posts a reply to a thread
func (c *Client) PostThreadReply(channelID, text, threadTimestamp string) error {
	// Slack API ドキュメントに基づく正しい実装
	_, _, err := c.api.PostMessage(channelID,
		slack.MsgOptionText(text, false),
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			ThreadTimestamp: threadTimestamp,
		}))
	if err != nil {
		return c.handleAPIError(err)
	}
	return nil
}

// PostThreadReplyByURL posts a reply to a thread using a thread URL
func (c *Client) PostThreadReplyByURL(text, threadURL string) error {
	// スレッドURLを解析
	threadInfo, err := ParseThreadURL(threadURL)
	if err != nil {
		return fmt.Errorf("スレッドURLの解析に失敗しました: %v", err)
	}

	// メッセージ情報を取得して正確なタイムスタンプを確認
	msg, err := c.GetMessageInfo(threadInfo.ChannelID, threadInfo.Timestamp)
	if err != nil {
		return fmt.Errorf("メッセージ情報の取得に失敗しました: %v", err)
	}

	// 正確なタイムスタンプでスレッド返信
	return c.PostThreadReply(threadInfo.ChannelID, text, msg.Timestamp)
}

// GetMessageInfo gets message information by timestamp
func (c *Client) GetMessageInfo(channelID, timestamp string) (*slack.Message, error) {
	// メッセージの詳細を取得するために、チャンネル履歴から検索
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     1000, // 十分な件数を取得
	}

	messages, err := c.api.GetConversationHistory(params)
	if err != nil {
		return nil, c.handleAPIError(err)
	}

	// 指定されたタイムスタンプのメッセージを検索
	// URLから取得したタイムスタンプの .000000 部分を削除して比較
	searchTimestamp := strings.TrimSuffix(timestamp, ".000000")
	for _, msg := range messages.Messages {
		if strings.HasPrefix(msg.Timestamp, searchTimestamp) {
			return &msg, nil
		}
	}

	return nil, fmt.Errorf("指定されたタイムスタンプのメッセージが見つかりませんでした: %s", timestamp)
}

// ReactionInfo contains information about a reaction
type ReactionInfo struct {
	Name  string      `json:"name"`
	Users []UserInfo  `json:"users"`
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
	_, err = c.GetMessageInfo(threadInfo.ChannelID, threadInfo.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("メッセージ情報の取得に失敗しました: %v", err)
	}

	// リアクション情報を取得
	itemRef := slack.ItemRef{
		Channel:   threadInfo.ChannelID,
		Timestamp: threadInfo.Timestamp,
	}
	
	reactions, err := c.api.GetReactions(itemRef, slack.GetReactionsParameters{})
	if err != nil {
		return nil, c.handleAPIError(err)
	}

	// リアクション情報を変換
	var reactionInfos []ReactionInfo
	for _, reaction := range reactions {
		// ユーザー情報を取得
		var users []UserInfo
		for _, userID := range reaction.Users {
			user, err := c.GetUserInfo(userID)
			if err != nil {
				// ユーザー情報が取得できない場合はスキップ
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

// TestConnection tests the Slack API connection
func (c *Client) TestConnection() error {
	// auth.test APIを呼び出して接続をテスト
	_, err := c.api.AuthTest()
	if err != nil {
		return c.handleAPIError(err)
	}

	return nil
}
