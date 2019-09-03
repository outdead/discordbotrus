package hook

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type Option func(client *Hook)

func SetSession(session *discordgo.Session) Option {
	return func(hook *Hook) {
		hook.session = session
	}
}

func SetFormatter(formatter logrus.Formatter) Option {
	return func(hook *Hook) {
		hook.formatter = formatter
	}
}

func SetLevels(levels []logrus.Level) Option {
	return func(hook *Hook) {
		hook.levels = levels
	}
}
