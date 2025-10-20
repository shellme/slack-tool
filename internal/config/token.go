package config

import (
	"fmt"
	"strings"
)

// MaskToken masks a token for display purposes
func MaskToken(token string) string {
	if len(token) <= 14 {
		// トークンが短すぎる場合は全てマスク
		return strings.Repeat("*", len(token))
	}

	// 最初の10文字 + "..." + 最後の4文字
	prefix := token[:10]
	suffix := token[len(token)-4:]

	return fmt.Sprintf("%s...%s", prefix, suffix)
}

// ValidateToken validates if the token format is correct
func ValidateToken(token string) error {
	if token == "" {
		return fmt.Errorf("トークンが空です")
	}

	if !strings.HasPrefix(token, "xoxp-") {
		return fmt.Errorf("User Tokenは 'xoxp-' で始まる必要があります")
	}

	if len(token) < 20 {
		return fmt.Errorf("トークンの長さが不正です")
	}

	return nil
}
