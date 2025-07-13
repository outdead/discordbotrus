package discordbotrus

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Option can be used to create a customized connection.
type Option func(client *Hook)

// WithSession sets discordgo session to Hook.
func WithSession(session *discordgo.Session) Option {
	return func(hook *Hook) {
		hook.session = session
	}
}

// WithFormatter sets custom formatter to Hook.
func WithFormatter(formatter logrus.Formatter) Option {
	return func(hook *Hook) {
		hook.formatter = formatter
	}
}

// WithLevels sets logrus levels to Hook.
func WithLevels(levels []logrus.Level) Option {
	return func(hook *Hook) {
		hook.levels = levels
	}
}
