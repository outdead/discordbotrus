package hook

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func TestEmbedFormatter_ErrorNotLost(t *testing.T) {
	formatter := &EmbedFormatter{}

	b, err := formatter.Format(logrus.WithField("error", errors.New("wild walrus")))
	if err != nil {
		t.Fatal("unable to format entry: ", err)
	}

	embed := discordgo.MessageEmbed{}
	if err = json.Unmarshal(b, &embed); err != nil {
		t.Fatal("unable to unmarshal formatted entry: ", err)
	}

	if len(embed.Fields) == 0 {
		t.Fatal("error field not set")
	}

	field := embed.Fields[0]
	if field.Name != "error" || field.Value != "wild walrus" {
		t.Fatal("error field not set")
	}
}

func TestEmbedFormatter_FieldClashWithTime(t *testing.T) {
	formatter := &EmbedFormatter{}

	b, err := formatter.Format(logrus.WithField("time", "right now!"))
	if err != nil {
		t.Fatal("unable to format entry: ", err)
	}

	embed := discordgo.MessageEmbed{}
	if err = json.Unmarshal(b, &embed); err != nil {
		t.Fatal("unable to unmarshal formatted entry: ", err)
	}

	if embed.Title != "0001-01-01 00:00:00 PANIC" {
		t.Fatal("title field not set to current time, was: ", embed.Title)
	}

	if len(embed.Fields) == 0 {
		t.Fatal("error field not set")
	}

	field := embed.Fields[0]
	if field.Name != "time" || field.Value != "right now!" {
		t.Fatal("fields.time not set to original time field")
	}
}

func TestEmbedFormatter_FieldClashWithMsg(t *testing.T) {
	formatter := &EmbedFormatter{}

	b, err := formatter.Format(logrus.WithField("msg", "something"))
	if err != nil {
		t.Fatal("unable to format entry: ", err)
	}

	embed := discordgo.MessageEmbed{}
	if err = json.Unmarshal(b, &embed); err != nil {
		t.Fatal("unable to unmarshal formatted entry: ", err)
	}

	if len(embed.Fields) == 0 {
		t.Fatal("error field not set")
	}

	field := embed.Fields[0]
	if field.Name != "msg" || field.Value != "something" {
		t.Fatal("fields.msg not set to original msg field")
	}
}

func TestEmbedFormatter_CustomSorting(t *testing.T) {
	formatter := &EmbedFormatter{
		SortingFunc: func(keys []string) {
			sort.Slice(keys, func(i, j int) bool {
				if keys[j] == "prefix" {
					return false
				}
				if keys[i] == "prefix" {
					return true
				}
				return strings.Compare(keys[i], keys[j]) == -1
			})
		},
	}

	entry := &logrus.Entry{
		Message: "Testing custom sort function",
		Time:    time.Now(),
		Level:   logrus.InfoLevel,
		Data: logrus.Fields{
			"test":      "testvalue",
			"prefix":    "the application prefix",
			"blablabla": "blablabla",
		},
	}
	b, err := formatter.Format(entry)
	if err != nil {
		t.Errorf("expected nil got error: %s", err)
	}

	embed := discordgo.MessageEmbed{}
	err = json.Unmarshal(b, &embed)
	if err != nil {
		t.Fatal("unable to unmarshal formatted entry: ", err)
	}

	if len(embed.Fields) == 0 {
		t.Fatal("error field not set")
	}

	field := embed.Fields[0]
	if field.Name != "prefix" || field.Value != "the application prefix" {
		t.Errorf("no expected fields order: %s", string(b))
	}
}

func TestEmbedFormatter_Limits(t *testing.T) {
	// Test embedMaxFieldCount
	func() {
		formatter := DefaultEmbedFormatter
		entry := &logrus.Entry{}

		for i := 0; i < embedMaxFieldCount+5; i++ {
			key := fmt.Sprintf("key %d", i+1)
			entry = entry.WithField(key, "value")
		}

		embed := formatter.Embed(entry)
		if len(embed.Fields) != embedMaxFieldCount {
			t.Errorf("unexpected fields count %d", len(embed.Fields))
		}
	}()

	// Test embedMaxDescriptionLen, embedMaxFieldNameLen and embedMaxFieldValueLen
	func() {
		formatter := DefaultEmbedFormatter
		entry := &logrus.Entry{}

		key := string(make([]byte, embedMaxFieldNameLen+5))
		value := string(make([]byte, embedMaxFieldValueLen+5))
		entry = entry.WithField(key, value)
		entry.Message = string(make([]byte, embedMaxDescriptionLen+5))

		embed := formatter.Embed(entry)

		if len(embed.Description) != embedMaxDescriptionLen {
			t.Errorf("unexpected embed description len %d", len(embed.Description))
			return
		}

		if len(embed.Fields) == 0 {
			t.Error("key field not set")
			return
		}

		field := embed.Fields[0]
		if len(field.Name) != embedMaxFieldNameLen {
			t.Errorf("unexpected key len %d", len(field.Name))
			return
		}

		if len(field.Value) != embedMaxFieldValueLen {
			t.Errorf("unexpected value len %d", len(field.Value))
			return
		}
	}()

	// Test empty embed field value.
	func() {
		formatter := DefaultEmbedFormatter

		entry := &logrus.Entry{}
		entry = entry.WithField("key", "")

		embed := formatter.Embed(entry)
		if len(embed.Fields) != 0 {
			t.Errorf("unexpected fields len %d", len(embed.Fields))
			return
		}
	}()
}
