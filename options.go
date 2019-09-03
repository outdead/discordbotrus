package hook

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Option can be used to a create a customized connection.
type Option func(client *Hook)

// SetSession sets discordgo session to Hook.
func SetSession(session *discordgo.Session) Option {
	return func(hook *Hook) {
		hook.session = session
	}
}

// SetFormatter sets custom formatter to Hook.
func SetFormatter(formatter logrus.Formatter) Option {
	return func(hook *Hook) {
		hook.formatter = formatter
	}
}

// SetLevels sets logrus levels to Hook.
func SetLevels(levels []logrus.Level) Option {
	return func(hook *Hook) {
		hook.levels = levels
	}
}
