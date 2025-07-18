package discordbotrus

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	// DiscordMaxMessageLen max discord message length.
	DiscordMaxMessageLen = 2000

	// DefaultTimeLayout default time layout to Formatter implementations.
	DefaultTimeLayout = "2006-01-02 15:04:05"
)

var (
	// ErrEmptyToken is returned when discord bot token is empty with enabled hook.
	ErrEmptyToken = errors.New("discord bot token is empty")

	// ErrEmptyChannelID is returned when discord channel id is empty with enabled hook.
	ErrEmptyChannelID = errors.New("discord channel id is empty")

	// ErrMessageTooLong is returned when message that has been sent to discord longer
	// than 2000 characters.
	ErrMessageTooLong = errors.New("discord message too long")
)

// Hook implements logrus.Hook and delivers logs to discord channel.
type Hook struct {
	config    *Config
	levels    []logrus.Level
	session   *discordgo.Session
	owner     bool
	formatter logrus.Formatter
}

// New creates new logrus hook for discord.
func New(cfg *Config, options ...Option) (*Hook, error) {
	hook := Hook{
		config: cfg,
	}

	for _, option := range options {
		option(&hook)
	}

	// Allow use hook if it is deactivated. Necessary in order to simplify compatibility.
	if cfg.Disabled {
		return &hook, nil
	}

	if err := cfg.Validate(); err != nil {
		return &hook, err
	}

	if hook.session == nil && cfg.Token == "" {
		return &hook, ErrEmptyToken
	}

	// Add missed properties if hook is enabled.
	if err := hook.setDefaults(); err != nil {
		return &hook, err
	}

	return &hook, nil
}

// Fire implements logrus.Hook.
func (hook *Hook) Fire(entry *logrus.Entry) error {
	// Do nothing if hook is disabled in config.
	if hook.config.Disabled {
		return nil
	}

	if hook.config.Format == EmbedFormatterCode || hook.config.Format == "" {
		embedFormatter, ok := hook.formatter.(*EmbedFormatter)
		if ok {
			msg := embedFormatter.Embed(entry)

			if _, err := hook.session.ChannelMessageSendEmbed(hook.config.ChannelID, msg); err != nil {
				return fmt.Errorf("discordbotrus: %w", err)
			}

			return nil
		}
	}

	raw, err := hook.formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("discordbotrus: %w", err)
	}

	if _, err := hook.session.ChannelMessageSend(hook.config.ChannelID, string(raw)); err != nil {
		return fmt.Errorf("discordbotrus: %w", err)
	}

	return nil
}

// Levels implements logrus.Hook.
func (hook *Hook) Levels() []logrus.Level {
	if hook.levels == nil {
		return logrus.AllLevels
	}

	return hook.levels
}

// Close implements io.Closer.
// Closes connection to discord if hook is owner of it.
func (hook *Hook) Close() error {
	// Do nothing if hook is disabled in config.
	if hook.config.Disabled {
		return nil
	}

	// Close discord connection only if it was opened during initialization.
	if hook.owner {
		if err := hook.session.Close(); err != nil {
			return fmt.Errorf("discordbotrus: %w", err)
		}
	}

	return nil
}

// setDefaults adds missed properties and sets default values to hook.
func (hook *Hook) setDefaults() error {
	if hook.levels == nil {
		var err error
		if hook.levels, err = ParseLevels(hook.config.Levels, hook.config.MinLevel); err != nil {
			return err
		}
	}

	if hook.session == nil {
		session, err := discordgo.New("Bot " + hook.config.Token)
		if err != nil {
			return fmt.Errorf("create session: %w", err)
		}

		if err := session.Open(); err != nil {
			return fmt.Errorf("open session: %w", err)
		}

		hook.session = session
		hook.owner = true
	}

	if hook.formatter == nil {
		switch hook.config.Format {
		case TextFormatterCode:
			hook.formatter = DefaultTextFormatter
		case JSONFormatterCode:
			hook.formatter = DefaultJSONFormatter
		case EmbedFormatterCode:
			hook.formatter = DefaultEmbedFormatter
		default:
			hook.formatter = DefaultEmbedFormatter
		}
	}

	return nil
}
