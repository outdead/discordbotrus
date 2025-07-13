package discordbotrus

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func TestNewHookDisabled(t *testing.T) {
	cfg := NewDefaultConfig("", "")
	cfg.Disabled = true

	hook, err := New(cfg)
	if err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}

	hook.Fire(logrus.WithField("test", "test"))

	if err := hook.Close(); err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}
}

func TestNewHookWithOptions(t *testing.T) {
	cfg := getConfig(EmbedFormatterCode)

	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}

	defer func() {
		if err := session.Close(); err != nil {
			t.Fatalf("expected nil got error: %s", err)
		}
	}()

	hook, err := New(
		cfg,
		WithSession(session),
		WithFormatter(DefaultJSONFormatter),
		WithLevels([]logrus.Level{logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.TraceLevel}),
	)
	if err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}

	if err := hook.Close(); err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}
}

func TestNewHookWithSessionOwned(t *testing.T) {
	cfg := getConfig("")
	hook, err := New(cfg)
	if err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}

	time.Sleep(10 * time.Millisecond)

	defer func() {
		if err := hook.Close(); err != nil {
			t.Fatalf("expected nil got error: %s", err)
		}
	}()
}

func TestNewHook_ErrEmptyChannel(t *testing.T) {
	cfg := NewDefaultConfig("", "")

	hook, err := New(cfg)
	if !errors.Is(err, ErrEmptyChannelID) {
		t.Fatalf("expected error: %s got error: %s", ErrEmptyChannelID, err)
	}

	if err := hook.Close(); err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}
}

func TestNewHook_ErrEmptyToken(t *testing.T) {
	cfg := NewDefaultConfig("", "1234567890")

	hook, err := New(cfg)
	if !errors.Is(err, ErrEmptyToken) {
		t.Fatalf("expected error: %s got error: %s", ErrEmptyToken, err)
	}

	if err := hook.Close(); err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}
}

func TestHook_Fire_Embed(t *testing.T) {
	testHookFire(t, getConfig(EmbedFormatterCode))
}

func TestHook_Fire_JSON(t *testing.T) {
	testHookFire(t, getConfig(JSONFormatterCode))
}

func TestHook_Fire_Text(t *testing.T) {
	testHookFire(t, getConfig(TextFormatterCode))
}

func testHookFire(t *testing.T, cfg *Config) {
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}

	hook, err := New(cfg, WithSession(session))
	if err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}

	defer func() {
		if err := session.Close(); err != nil {
			t.Errorf("expected nil got error: %s", err)
		}

		if err := hook.Close(); err != nil {
			t.Fatalf("expected nil got error: %s", err)
		}
	}()

	var (
		textInfo = "1 The Tavern [11568x8230](https://map.projectzomboid.com/#11568x8230) - [11636x8291](https://map.projectzomboid.com/#11636x8291)\n" +
			"2 Used cars [10057x13114](https://map.projectzomboid.com/#10057x13114) - [10100x13149](https://map.projectzomboid.com/#10100x13149)\n" +
			"3 Flame Bar [5787x5389](https://map.projectzomboid.com/#5787x5389) - [5862x5440](https://map.projectzomboid.com/#5862x5440)\n" +
			"4 Admin House [10202x12772](https://map.projectzomboid.com/#10202x12772) - [10232x12820](https://map.projectzomboid.com/#10232x12820)\n" +
			"5 Hospital [7369x8314](https://map.projectzomboid.com/#7369x8314) - [7446x8422](https://map.projectzomboid.com/#7446x8422)"

		textWarning = "[31-07-19 17:57:43.823] 76561190000000000 \"Username\" removed MetalWallLvl2 (constructedobjects_01_49) at 11604,8240,0."

		textError = "connected time is not set in user storage"

		fieldsError = map[string]interface{}{
			"component":   "testhook",
			"coordinates": "11612,8250,0",
			"event_time":  "2019-07-29T21:02:44Z",
			"owner_id":    "",
			"steam_id":    "76561190000000000",
			"username":    "Username",
			"action":      "disconnected player",
			"line":        "[29-07-19 21:02:44.504] 76561190000000000 \"Username\" disconnected player (11612,8250,0)",
		}
	)

	if err := hook.Fire(&logrus.Entry{Message: textInfo, Level: logrus.InfoLevel, Time: time.Now()}); err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}

	if err := hook.Fire(&logrus.Entry{Message: textWarning, Level: logrus.WarnLevel, Time: time.Now()}); err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}

	if err := hook.Fire(&logrus.Entry{Message: textError, Level: logrus.ErrorLevel, Time: time.Now(), Data: fieldsError}); err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}
}

func getConfig(format string) *Config {
	token := os.Getenv("LDH_TOKEN")
	channelID := os.Getenv("LDH_CHANNEL")
	cfg := NewDefaultConfig(token, channelID)
	cfg.Format = format

	return cfg
}
