package hook

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Colors for log levels.
const (
	ColorGreen  = 0x0008000
	ColorYellow = 0xffaa00
	ColorRed    = 0xff0000
)

const levelsCount = 7

// ParseLevels parses logging levels from the config.
func ParseLevels(lvs []string, minLvl string) ([]logrus.Level, error) {
	levels := make([]logrus.Level, 0, levelsCount)

	if minLvl == "" {
		all := logrus.AllLevels
		minLvl = all[len(all)-1].String()
	}

	minLevel, err := logrus.ParseLevel(minLvl)
	if err != nil {
		return levels, fmt.Errorf("discordbotrus: %w", err)
	}

	if len(lvs) != 0 {
		for _, lvl := range lvs {
			level, err := logrus.ParseLevel(lvl)
			if err != nil {
				return levels, fmt.Errorf("discordbotrus: %w", err)
			}

			if minLevel >= level {
				levels = append(levels, level)
			}
		}

		return levels, nil
	}

	for _, level := range logrus.AllLevels {
		if minLevel >= level {
			levels = append(levels, level)
		}
	}

	return levels, nil
}

// LevelColor returns the respective color for the logrus level.
func LevelColor(l logrus.Level) int {
	switch l {
	case logrus.PanicLevel:
		return ColorRed
	case logrus.FatalLevel:
		return ColorRed
	case logrus.ErrorLevel:
		return ColorRed
	case logrus.WarnLevel:
		return ColorYellow
	case logrus.InfoLevel:
		return ColorGreen
	case logrus.DebugLevel:
		return ColorGreen
	case logrus.TraceLevel:
		return ColorGreen
	default:
		return 0
	}
}
