package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	SlackToken string `json:"slack_token"`
}

// ConfigManager handles configuration file operations
type ConfigManager struct {
	configPath string
}

// NewConfigManager creates a new ConfigManager instance
func NewConfigManager() *ConfigManager {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("ホームディレクトリを取得できませんでした: %v", err))
	}

	configDir := filepath.Join(homeDir, ".config", "slack-tool")
	configPath := filepath.Join(configDir, "config.json")

	return &ConfigManager{
		configPath: configPath,
	}
}

// Load loads configuration from file
func (cm *ConfigManager) Load() (*Config, error) {
	// 設定ファイルが存在しない場合は空の設定を返す
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		return &Config{}, nil
	}

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("設定ファイルの解析に失敗しました: %v", err)
	}

	return &config, nil
}

// Save saves configuration to file
func (cm *ConfigManager) Save(config *Config) error {
	// 設定ディレクトリを作成
	configDir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("設定ディレクトリの作成に失敗しました: %v", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("設定のJSON変換に失敗しました: %v", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0600); err != nil {
		return fmt.Errorf("設定ファイルの保存に失敗しました: %v", err)
	}

	return nil
}

// GetConfigPath returns the configuration file path
func (cm *ConfigManager) GetConfigPath() string {
	return cm.configPath
}
