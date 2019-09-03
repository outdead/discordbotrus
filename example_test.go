package hook_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/outdead/discordbotrus"
	"github.com/sirupsen/logrus"
)

func TestExample(t *testing.T) {
	cfg := hook.NewDefaultConfig(os.Getenv("LDH_TOKEN"), os.Getenv("LDH_CHANNEL"))
	hooker, err := hook.New(cfg)
	if err != nil {
		t.Fatalf("expected nil got error: %s", err)
	}

	logger := &logrus.Logger{Out: ioutil.Discard, Formatter: new(logrus.JSONFormatter), Hooks: make(logrus.LevelHooks), Level: logrus.InfoLevel}
	logger.AddHook(hooker)

	defer func() {
		if err := hooker.Close(); err != nil {
			t.Fatalf("expected nil got error: %s", err)
		}
	}()

	logger.Info("My spoon is too big")
}
