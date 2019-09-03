package hook_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/outdead/discordbotrus"
	"github.com/sirupsen/logrus"
)

func TestExample(t *testing.T) {
	cfg := hook.NewDefaultConfig(os.Getenv("LDH_TOKEN"), os.Getenv("LDH_CHANNEL"))
	hooker, err := hook.New(cfg)
	if err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}

	defer func() {
		if err := hooker.Close(); err != nil {
			t.Fatalf("expected nil got error: %s", err)
		}
	}()

	logger := &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
	logger.AddHook(hooker)

	logger.Info("My spoon is too big")
}

func TestExampleWithSession(t *testing.T) {
	token := os.Getenv("LDH_TOKEN")
	channelID := os.Getenv("LDH_CHANNEL")

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}

	// In this case, you can use the session without opening a web socket.
	// But to establish a stable connection, it is better to do this.
	if err := session.Open(); err != nil {
		t.Fatalf("open discord session error: %s", err)
	}

	defer func() {
		if err != session.Close() {
			t.Fatalf("expected nil got error: %s", err)
		}
	}()

	hooker, err := hook.New(
		&hook.Config{ChannelID: channelID},
		hook.SetSession(session),
	)
	if err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}

	logger := &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
	logger.AddHook(hooker)

	logger.Info("My spoon is too big")
}
