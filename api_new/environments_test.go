package api_new

import (
	"testing"
)

func TestGetSettings(t *testing.T) {
	locale := getLocale(`2020.2`)
	if locale != `en` {
		t.Fatalf(`expected 'en' but got '%v'`, locale)
	}
}
