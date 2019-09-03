package hook

import "testing"

func TestConfig_Validate(t *testing.T) {
	cfg := NewDefaultConfig("", "")
	if err := cfg.Validate(); err != ErrEmptyChannelID {
		t.Fatalf("expect err: %s got %s", ErrEmptyChannelID, err)
	}

	cfg.Disabled = true
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expect err: nil got %s", err)
	}

	cfg = NewDefaultConfig("", "1234567890")
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expect err: nil got %s", err)
	}
}
