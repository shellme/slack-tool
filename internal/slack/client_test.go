package slack

import "testing"

func TestClientParseTimestamp(t *testing.T) {
	c := &Client{}

	tests := []struct {
		input string
		want  string
	}{
		{input: "1700000000", want: "1700000000"},
		{input: "2026-01-01T00:00:00Z", want: "1767225600"},
		{input: "2026-01-01T09:00:00+09:00", want: "1767225600"},
	}

	for _, tt := range tests {
		got, err := c.parseTimestamp(tt.input)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("%s: got %s, want %s", tt.input, got, tt.want)
		}
	}

	if _, err := c.parseTimestamp("not-a-date"); err == nil {
		t.Error("expected error for invalid date")
	}
}
