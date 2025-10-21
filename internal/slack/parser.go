package slack

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ThreadInfo contains parsed information from a Slack thread URL
type ThreadInfo struct {
	ChannelID string
	Timestamp string
}

// ChannelInfo contains parsed information from a Slack channel URL
type ChannelInfo struct {
	ChannelID string
}

// ThreadURLInfo contains parsed information from a Slack thread URL
type ThreadURLInfo struct {
	ChannelID       string
	Timestamp       string
	ThreadTimestamp string // スレッド返信URLの場合は thread_ts パラメータ
}

// ParseSlackURL parses a Slack thread URL and extracts channel ID and timestamp
func ParseSlackURL(url string) (*ThreadInfo, error) {
	// Slack URLの正規表現パターン
	// 例: https://your-workspace.slack.com/archives/C12345678/p1234567890123456
	pattern := `https://[^/]+\.slack\.com/archives/([A-Z0-9]+)/p(\d+)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("正規表現のコンパイルに失敗しました: %v", err)
	}

	matches := re.FindStringSubmatch(url)
	if len(matches) != 3 {
		return nil, fmt.Errorf("無効なSlackスレッドURLです。正しい形式: https://your-workspace.slack.com/archives/C12345678/p1234567890123456")
	}

	channelID := matches[1]
	timestampStr := matches[2]

	// タイムスタンプの形式を変換
	// p1234567890123456 -> 1234567890.123456
	if len(timestampStr) < 10 {
		return nil, fmt.Errorf("無効なタイムスタンプです: %s", timestampStr)
	}

	// 秒部分とマイクロ秒部分に分割
	seconds := timestampStr[:10]
	microseconds := timestampStr[10:]

	// マイクロ秒部分を6桁に調整（SlackのAPIが期待する形式）
	if len(microseconds) > 6 {
		microseconds = microseconds[:6]
	} else {
		// 6桁にパディング
		for len(microseconds) < 6 {
			microseconds += "0"
		}
	}

	// 数値として有効かチェック
	if _, err := strconv.ParseInt(seconds, 10, 64); err != nil {
		return nil, fmt.Errorf("無効なタイムスタンプの秒部分です: %s", seconds)
	}

	timestamp := seconds + "." + microseconds

	return &ThreadInfo{
		ChannelID: channelID,
		Timestamp: timestamp,
	}, nil
}

// ParseChannelURL parses a Slack channel URL and extracts channel ID
func ParseChannelURL(url string) (*ChannelInfo, error) {
	// SlackチャンネルURLの正規表現パターン
	// 例: https://your-workspace.slack.com/archives/C12345678
	pattern := `https://[^/]+\.slack\.com/archives/([A-Z0-9]+)/?$`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("正規表現のコンパイルに失敗しました: %v", err)
	}

	matches := re.FindStringSubmatch(url)
	if len(matches) != 2 {
		return nil, fmt.Errorf("無効なSlackチャンネルURLです。正しい形式: https://your-workspace.slack.com/archives/C12345678")
	}

	channelID := matches[1]

	return &ChannelInfo{
		ChannelID: channelID,
	}, nil
}

// ValidateChannelID validates if the channel ID format is correct
func ValidateChannelID(channelID string) error {
	if channelID == "" {
		return fmt.Errorf("チャンネルIDが空です")
	}

	// チャンネルIDの形式をチェック（C, G, Dで始まる）
	if !strings.HasPrefix(channelID, "C") && !strings.HasPrefix(channelID, "G") && !strings.HasPrefix(channelID, "D") {
		return fmt.Errorf("無効なチャンネルID形式です: %s", channelID)
	}

	return nil
}

// ParseThreadURL parses a Slack thread URL and extracts channel ID and timestamp
func ParseThreadURL(url string) (*ThreadURLInfo, error) {
	// スレッド返信URLの場合は thread_ts パラメータを抽出
	if strings.Contains(url, "thread_ts=") {
		return parseThreadReplyURL(url)
	}

	// 通常のスレッドURLの解析
	// SlackスレッドURLの正規表現パターン
	// 例: https://your-workspace.slack.com/archives/C12345678/p1234567890123456
	pattern := `https://[^/]+\.slack\.com/archives/([A-Z0-9]+)/p(\d+)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("正規表現のコンパイルに失敗しました: %v", err)
	}

	matches := re.FindStringSubmatch(url)
	if len(matches) != 3 {
		return nil, fmt.Errorf("無効なSlackスレッドURLです。正しい形式: https://your-workspace.slack.com/archives/C12345678/p1234567890123456")
	}

	channelID := matches[1]
	timestampStr := matches[2]

	// タイムスタンプの形式を変換
	// p1760786585959009 -> 1760786585.959009
	if len(timestampStr) < 10 {
		return nil, fmt.Errorf("無効なタイムスタンプです: %s", timestampStr)
	}

	// 秒部分とマイクロ秒部分に分割
	seconds := timestampStr[:10]
	microseconds := timestampStr[10:]

	// マイクロ秒部分を6桁に調整（SlackのAPIが期待する形式）
	if len(microseconds) > 6 {
		microseconds = microseconds[:6]
	} else {
		// 6桁にパディング
		for len(microseconds) < 6 {
			microseconds += "0"
		}
	}

	// 数値として有効かチェック
	if _, err := strconv.ParseInt(seconds, 10, 64); err != nil {
		return nil, fmt.Errorf("無効なタイムスタンプの秒部分です: %s", seconds)
	}

	timestamp := seconds + "." + microseconds

	return &ThreadURLInfo{
		ChannelID: channelID,
		Timestamp: timestamp,
	}, nil
}

// parseThreadReplyURL parses a Slack thread reply URL and extracts channel ID and thread timestamp
func parseThreadReplyURL(url string) (*ThreadURLInfo, error) {
	// スレッド返信URLの正規表現パターン
	// 例: https://your-workspace.slack.com/archives/C12345678/p1234567890123456?thread_ts=1760786585.959009&cid=C12345678
	pattern := `https://[^/]+\.slack\.com/archives/([A-Z0-9]+)/p(\d+).*[?&]thread_ts=([0-9.]+)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("正規表現のコンパイルに失敗しました: %v", err)
	}

	matches := re.FindStringSubmatch(url)
	if len(matches) != 4 {
		return nil, fmt.Errorf("無効なSlackスレッド返信URLです。正しい形式: https://your-workspace.slack.com/archives/C12345678/p1234567890123456?thread_ts=1760786585.959009&cid=C12345678")
	}

	channelID := matches[1]
	messageTimestampStr := matches[2]
	threadTimestamp := matches[3]

	// メッセージのタイムスタンプを変換
	// p1760949917358289 -> 1760949917.358289
	if len(messageTimestampStr) < 10 {
		return nil, fmt.Errorf("無効なメッセージタイムスタンプです: %s", messageTimestampStr)
	}

	// 秒部分とマイクロ秒部分に分割
	seconds := messageTimestampStr[:10]
	microseconds := messageTimestampStr[10:]

	// マイクロ秒部分を6桁に調整
	if len(microseconds) > 6 {
		microseconds = microseconds[:6]
	} else {
		// 6桁にパディング
		for len(microseconds) < 6 {
			microseconds += "0"
		}
	}

	// 数値として有効かチェック
	if _, err := strconv.ParseInt(seconds, 10, 64); err != nil {
		return nil, fmt.Errorf("無効なメッセージタイムスタンプの秒部分です: %s", seconds)
	}

	messageTimestamp := seconds + "." + microseconds

	return &ThreadURLInfo{
		ChannelID:       channelID,
		Timestamp:       messageTimestamp, // 実際のメッセージのタイムスタンプ
		ThreadTimestamp: threadTimestamp,  // スレッドのタイムスタンプ
	}, nil
}
