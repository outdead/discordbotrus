package hook

import "github.com/sirupsen/logrus"

// TextFormatterCode formatter code to identify from config.
const TextFormatterCode = "text"

// DefaultTextFormatter used as default TextFormatter.
var DefaultTextFormatter = &TextFormatter{
	TextFormatter: logrus.TextFormatter{
		TimestampFormat:  DefaultTimeLayout,
		QuoteEmptyFields: true,
		//DisableSorting:   true,
	},
	Quoted: true,
}

// TextFormatter formats logs into text.
type TextFormatter struct {
	logrus.TextFormatter

	// Quoted will quote string to discord tag text.
	Quoted bool
}

// Format renders a single log entry.
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data, err := f.TextFormatter.Format(entry)
	if err != nil {
		return data, err
	}

	if f.Quoted {
		data = []byte("```text\n" + string(data) + "```")
	}

	if len([]rune(string(data))) > DiscordMaxMessageLen {
		return data, ErrMessageTooLong
	}

	return data, err
}
