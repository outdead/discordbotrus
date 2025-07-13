package discordbotrus

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestTextFormatter_FormatQuoted(t *testing.T) {
	formatter := DefaultTextFormatter

	b, err := formatter.Format(logrus.WithField("key", "wild walrus"))
	if err != nil {
		t.Fatalf("unable to format entry: %s", err)
	}

	bstr := string(b)
	if strings.Index(bstr, "```text\n") != 0 || strings.Index(bstr, "\n```") != len(bstr)-4 {
		t.Fatal("text tag not set")
	}

	if !bytes.Contains(b, []byte("key=\"wild walrus\"")) {
		t.Fatal("key field not set")
	}
}

func TestTextFormatter_FormatLongMessage(t *testing.T) {
	formatter := &TextFormatter{}

	entry := &logrus.Entry{}
	entry = entry.WithField("key", "value")
	entry.Message = string(make([]byte, DiscordMaxMessageLen))

	_, err := formatter.Format(entry)
	if err == nil {
		t.Fatal("expect error got nil")
	}

	if err != ErrMessageTooLong {
		t.Fatalf("unexpected error: %s got error: %s", err, ErrMessageTooLong)
	}
}
