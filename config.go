package hook

// Config handles discord bot connection configuration and logrus levels.
type Config struct {
	// Disabled can disable hook form configuration file.
	Disabled bool `json:"disabled" yaml:"disabled"`

	// Token is bot token from discord developers applications.
	Token string `json:"token" yaml:"token"`

	// ChannelID is id of discord channel to log hooks.
	ChannelID string `json:"channel_id" yaml:"channel_id"`

	// Format specifies formatter to discord message.
	// Supported formats: text, json, embed.
	Format string `json:"format" yaml:"format"`

	// MinLevel is the minimum priority level to enable logging.
	MinLevel string `json:"min_level" yaml:"min_level"`

	// Levels is a list of levels to enable logging. Intersects with MinLevel.
	Levels []string `json:"levels" yaml:"levels"`
}

// NewDefaultConfig returns default configuration for hook.
func NewDefaultConfig(token string, channelID string) *Config {
	return &Config{
		Disabled:  false,
		Token:     token,
		ChannelID: channelID,
		MinLevel:  "info",
		Format:    EmbedFormatterCode,
		Levels: []string{
			"error",
			"warning",
			"info",
			"trace",
		},
	}
}

// Validate checks config for required fields.
func (cfg *Config) Validate() error {
	// Do not validate disabled hook.
	if cfg.Disabled {
		return nil
	}

	if cfg.ChannelID == "" {
		return ErrEmptyChannelID
	}

	return nil
}
