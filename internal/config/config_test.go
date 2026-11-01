package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigManagerSaveLoad(t *testing.T) {
	cm := &ConfigManager{configPath: filepath.Join(t.TempDir(), "nested", "config.json")}

	cfg, err := cm.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SlackToken != "" {
		t.Errorf("expected empty config, got %+v", cfg)
	}

	if err := cm.Save(&Config{SlackToken: "dummy"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(cm.configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("got permission %o, want 600", info.Mode().Perm())
	}

	cfg, err = cm.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SlackToken != "dummy" {
		t.Errorf("got %s, want dummy", cfg.SlackToken)
	}
}

func TestValidateToken(t *testing.T) {
	if err := ValidateToken("xoxp-" + "0123456789abcdef"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	for _, token := range []string{"", "xoxb-0123456789abcdef", "xoxp-short"} {
		if err := ValidateToken(token); err == nil {
			t.Errorf("%q: expected error", token)
		}
	}
}

func TestMaskToken(t *testing.T) {
	if got := MaskToken("short"); got != "*****" {
		t.Errorf("got %s", got)
	}
	if got := MaskToken("abcdefghij-middle-wxyz"); got != "abcdefghij...wxyz" {
		t.Errorf("got %s", got)
	}
}
