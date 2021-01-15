package hook

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// JSONFormatterCode formatter code to identify from config.
const JSONFormatterCode = "json"

// DefaultJSONFormatter used as default JSONFormatter.
var DefaultJSONFormatter = &JSONFormatter{
	JSONFormatter: logrus.JSONFormatter{
		TimestampFormat: DefaultTimeLayout,
	},
	Quoted: true,
}

// JSONFormatter formats logs into parsable json.
type JSONFormatter struct {
	logrus.JSONFormatter

	// Quoted will quote string to discord tag json.
	Quoted bool
}

// Format renders a single log entry.
func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data, err := f.JSONFormatter.Format(entry)
	if err != nil {
		return data, fmt.Errorf("discordbotrus: %w", err)
	}

	if f.Quoted {
		data = []byte("```json\n" + string(data) + "```")
	}

	if len([]rune(string(data))) > DiscordMaxMessageLen {
		return data, ErrMessageTooLong
	}

	return data, nil
}
