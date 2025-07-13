package discordbotrus

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestJSONFormatter_FormatQuoted(t *testing.T) {
	formatter := DefaultJSONFormatter

	b, err := formatter.Format(logrus.WithField("key", "wild walrus"))
	if err != nil {
		t.Fatalf("unable to format entry: %s", err)
	}

	bstr := string(b)
	if strings.Index(bstr, "```json\n") != 0 || strings.Index(bstr, "\n```") != len(bstr)-4 {
		t.Fatal("json tag not set")
	}

	bstr = strings.TrimLeft(bstr, "```json")
	bstr = strings.TrimRight(bstr, "```")
	b = []byte(bstr)

	entry := make(map[string]interface{})
	if err = json.Unmarshal(b, &entry); err != nil {
		t.Fatalf("unable to unmarshal formatted entry: %s", err)
	}

	if entry["key"] != "wild walrus" {
		t.Fatal("key field not set")
	}
}

func TestJSONFormatter_FormatLongMessage(t *testing.T) {
	formatter := &JSONFormatter{}

	entry := &logrus.Entry{}
	entry = entry.WithField("key", "value")
	entry.Message = string(make([]byte, DiscordMaxMessageLen))

	_, err := formatter.Format(entry)
	if err == nil {
		t.Fatal("expect error got nil")
	}

	if !errors.Is(err, ErrMessageTooLong) {
		t.Fatalf("unexpected error: %s got error: %s", err, ErrMessageTooLong)
	}
}
