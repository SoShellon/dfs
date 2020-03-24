package common

import "testing"

func TestTokenize(t *testing.T) {
	parts := Tokenize("/")
	for range parts {
		t.Error("should be no parts")
	}

	parts = Tokenize("//a//b")
	if len(parts) != 2 {
		t.Errorf("should ignore //:%+v", parts)
	}
}
