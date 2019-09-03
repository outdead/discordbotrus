package hook

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	embedMaxFieldCount     = 25
	embedMaxFieldNameLen   = 256
	embedMaxFieldValueLen  = 1024
	embedMaxDescriptionLen = 2048

	// EmbedFormatterCode formatter code to identify from config.
	EmbedFormatterCode = "embed"
)

// DefaultEmbedFormatter used as default EmbedFormatter.
var DefaultEmbedFormatter = &EmbedFormatter{
	Inline:          true,
	TimestampFormat: DefaultTimeLayout,
	//DisableSorting:  true,
}

// EmbedFormatter formats logs into parsable json.
type EmbedFormatter struct {
	// Inline causes fields to be displayed one per line
	// as opposed to being inline.
	Inline bool

	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output.
	DisableTimestamp bool

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)
}

// Format renders a single log entry.
func (f *EmbedFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return json.Marshal(f.Embed(entry))
}

// Embed creates discord embed message from logrus entry.
func (f *EmbedFormatter) Embed(entry *logrus.Entry) *discordgo.MessageEmbed {
	title := strings.ToUpper(entry.Level.String())
	if !f.DisableTimestamp {
		timestampFormat := f.TimestampFormat
		if timestampFormat == "" {
			timestampFormat = DefaultTimeLayout
		}

		title = entry.Time.Format(timestampFormat) + " " + title
	}

	message := discordgo.MessageEmbed{
		Title: title,
		Color: LevelColor(entry.Level),
	}

	// Truncate message if it is too long.
	if len([]rune(entry.Message)) > embedMaxDescriptionLen {
		entry.Message = string([]rune(entry.Message)[:embedMaxDescriptionLen])
	}
	message.Description = entry.Message

	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		if f.SortingFunc == nil {
			sort.Strings(keys)
		} else {
			f.SortingFunc(keys)
		}
	}

	// Add fields to embed.
	counter := 0
	fields := make([]*discordgo.MessageEmbedField, 0, len(entry.Data))
	//for name, value := range entry.Data {
	for _, name := range keys {
		value := entry.Data[name]

		// Ensure that the maximum field number is not exceeded.
		if counter >= embedMaxFieldCount {
			break
		}

		// Make value a string.
		valueStr := fmt.Sprintf("%v", value)
		if len(valueStr) == 0 {
			// Fix for discordgo bug with empty field value. Discord responses
			// `HTTP 400 Bad Request, {"embed": ["fields"]}` and nothing is clear.
			// Therefore must skip the value or add a fake value.
			continue
		}

		// Truncate names and values which are too long.
		if len([]rune(name)) > embedMaxFieldNameLen {
			name = string([]rune(name)[:embedMaxFieldNameLen])
		}

		if len([]rune(valueStr)) > embedMaxFieldValueLen {
			valueStr = string([]rune(valueStr)[:embedMaxFieldValueLen])
		}

		var field = discordgo.MessageEmbedField{
			Name:   name,
			Value:  valueStr,
			Inline: f.Inline,
		}

		counter++
		fields = append(fields, &field)
	}
	message.Fields = fields

	return &message
}
