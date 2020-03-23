package common

import "testing"

func TestTokenize(t *testing.T) {
	parts := Tokenize("/")
	for range parts {
		t.Error("should be no parts")
	}
}
