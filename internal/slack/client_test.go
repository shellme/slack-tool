package slack

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	slackgo "github.com/slack-go/slack"
)

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

func TestClientParseTimestampUsesLocalTime(t *testing.T) {
	original := time.Local
	time.Local = time.FixedZone("JST", 9*60*60)
	defer func() { time.Local = original }()

	c := &Client{}
	got, err := c.parseTimestamp("2026-01-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "1767193200"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

// newTestClient starts a fake Slack API and returns a client connected to it
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return &Client{api: slackgo.New("dummy", slackgo.OptionAPIURL(server.URL+"/"))}
}

func TestGetThreadRepliesFollowsPagination(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("cursor") == "" {
			fmt.Fprint(w, `{"ok":true,"messages":[{"ts":"1.000001"}],"has_more":true,"response_metadata":{"next_cursor":"next"}}`)
			return
		}
		fmt.Fprint(w, `{"ok":true,"messages":[{"ts":"1.000002"}],"has_more":false}`)
	})

	messages, err := c.GetThreadReplies("C1", "1.000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(messages) != 2 {
		t.Errorf("got %d messages, want 2", len(messages))
	}
}

func TestFetchHistoryStopsAtLimit(t *testing.T) {
	var requests int
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		r.ParseForm()
		count, _ := strconv.Atoi(r.Form.Get("limit"))
		messages := strings.TrimSuffix(strings.Repeat(`{"ts":"1.0"},`, count), ",")
		fmt.Fprintf(w, `{"ok":true,"messages":[%s],"has_more":true,"response_metadata":{"next_cursor":"next"}}`, messages)
	})

	messages, err := c.GetChannelHistory("C1", 250)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(messages) != 250 || requests != 2 {
		t.Errorf("got %d messages in %d requests", len(messages), requests)
	}
}

func TestRetriesOnRateLimit(t *testing.T) {
	var requests int
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		fmt.Fprint(w, `{"ok":true,"messages":[{"ts":"1.000001"}],"has_more":false}`)
	})

	if _, err := c.GetThreadReplies("C1", "1.000001"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if requests != 2 {
		t.Errorf("got %d requests, want 2", requests)
	}
}

func TestGetMessageInfo(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if strings.HasSuffix(r.URL.Path, "conversations.replies") {
			fmt.Fprint(w, `{"ok":true,"messages":[{"ts":"1.000001"},{"ts":"1.000002","thread_ts":"1.000001"}],"has_more":false}`)
			return
		}
		if r.Form.Get("oldest") != r.Form.Get("latest") || r.Form.Get("inclusive") != "1" {
			t.Errorf("unexpected params: %v", r.Form)
		}
		fmt.Fprint(w, `{"ok":true,"messages":[{"ts":"1.000001"}],"has_more":false}`)
	})

	if msg, err := c.GetMessageInfo("C1", "1.000001", ""); err != nil || msg.Timestamp != "1.000001" {
		t.Errorf("got %+v, %v", msg, err)
	}
	if msg, err := c.GetMessageInfo("C1", "1.000002", "1.000001"); err != nil || msg.ThreadTimestamp != "1.000001" {
		t.Errorf("got %+v, %v", msg, err)
	}
	if _, err := c.GetMessageInfo("C1", "1.000003", ""); err == nil {
		t.Error("expected error for missing message")
	}
}

func TestChannelHistoryDeduplicatesBroadcastReplies(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "conversations.replies") {
			fmt.Fprint(w, `{"ok":true,"messages":[{"ts":"1.000001","thread_ts":"1.000001"},{"ts":"1.000002","thread_ts":"1.000001"}],"has_more":false}`)
			return
		}
		fmt.Fprint(w, `{"ok":true,"messages":[{"ts":"1.000002","thread_ts":"1.000001"},{"ts":"1.000001","thread_ts":"1.000001"}],"has_more":false}`)
	})

	messages, err := c.GetChannelHistoryWithThreads("C1", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(messages) != 2 {
		t.Errorf("got %d messages, want 2", len(messages))
	}
}
