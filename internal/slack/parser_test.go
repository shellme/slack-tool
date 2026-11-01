package slack

import "testing"

func TestParseSlackURL(t *testing.T) {
	info, err := ParseSlackURL("https://example.slack.com/archives/C12345678/p1234567890123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ChannelID != "C12345678" || info.Timestamp != "1234567890.123456" {
		t.Errorf("unexpected result: %+v", info)
	}

	if _, err := ParseSlackURL("https://example.com/foo"); err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestParseChannelURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{name: "basic", url: "https://example.slack.com/archives/C12345678", want: "C12345678"},
		{name: "trailing slash", url: "https://example.slack.com/archives/C12345678/", want: "C12345678"},
		{name: "message URL", url: "https://example.slack.com/archives/C12345678/p1234567890123456", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ParseChannelURL(tt.url)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if info.ChannelID != tt.want {
				t.Errorf("got %s, want %s", info.ChannelID, tt.want)
			}
		})
	}
}

func TestParseThreadURL(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantTS     string
		wantThread string
	}{
		{
			name:   "message URL",
			url:    "https://example.slack.com/archives/C12345678/p1760786585959009",
			wantTS: "1760786585.959009",
		},
		{
			name:       "reply URL",
			url:        "https://example.slack.com/archives/C12345678/p1760949917358289?thread_ts=1760786585.959009&cid=C12345678",
			wantTS:     "1760949917.358289",
			wantThread: "1760786585.959009",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ParseThreadURL(tt.url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if info.ChannelID != "C12345678" || info.Timestamp != tt.wantTS || info.ThreadTimestamp != tt.wantThread {
				t.Errorf("unexpected result: %+v", info)
			}
		})
	}
}

func TestValidateChannelID(t *testing.T) {
	for _, id := range []string{"C123", "G123", "D123"} {
		if err := ValidateChannelID(id); err != nil {
			t.Errorf("%s: unexpected error: %v", id, err)
		}
	}
	for _, id := range []string{"", "X123"} {
		if err := ValidateChannelID(id); err == nil {
			t.Errorf("%q: expected error", id)
		}
	}
}
