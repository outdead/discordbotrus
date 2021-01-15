package hook

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestParseLevels(t *testing.T) {
	// Test not empty lvs and not empty minLvl
	testParseLevels(
		t,
		[]logrus.Level{logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel},
		[]string{"error", "warning", "info", "trace"},
		"info",
	)

	// Test empty lvs and empty minLvl
	testParseLevels(t, logrus.AllLevels, []string{}, "")

	// Test not empty lvs and empty minLvl
	testParseLevels(
		t,
		[]logrus.Level{logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.TraceLevel},
		[]string{"error", "warning", "info", "trace"},
		"",
	)

	// Test parse min level error
	func() {
		expErrStr := `discordbotrus: not a valid logrus Level: "fake"`
		_, err := ParseLevels(nil, "fake")

		if err == nil {
			t.Errorf("expect error: %s got error: nil", expErrStr)
			return
		}

		if err.Error() != expErrStr {
			t.Errorf("expect error: %s got error: %s", expErrStr, err)
		}
	}()

	// Test parse levels error
	func() {
		expErrStr := `discordbotrus: not a valid logrus Level: "wrong"`
		_, err := ParseLevels([]string{"error", "warning", "wrong", "trace"}, "")

		if err == nil {
			t.Errorf("expect error: %s got error: nil", expErrStr)
			return
		}

		if err.Error() != expErrStr {
			t.Errorf("expect error: %s got error: %s", expErrStr, err)
		}
	}()
}

func TestLevelColor(t *testing.T) {
	levels := logrus.AllLevels

	for _, level := range levels {
		color := LevelColor(level)
		if color == 0 {
			t.Errorf("enexpected color for level: %v", level)
			break
		}
	}

	color := LevelColor(logrus.Level(999))
	if color != 0 {
		t.Error("enexpected color for nonexistent level")
	}
}

func testParseLevels(t *testing.T, exp []logrus.Level, lvs []string, minLvl string) {
	levels, err := ParseLevels(lvs, minLvl)
	if err != nil {
		t.Errorf("exp error: nil got error: %s", err)
		return
	}

	if len(levels) != len(exp) {
		fmt.Println(levels)
		t.Error("levels len not equal exp len")
		return
	}

	for k, v := range exp {
		if levels[k] != v {
			t.Errorf("exp %v got %v", v, levels[k])
			return
		}
	}
}
