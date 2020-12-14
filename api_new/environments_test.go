package api_new

import (
	"testing"
)

// These unit tests require Alteryx to be installed.  Change the version below with your version of ALteryx.

const version = `2020.2`

func TestGetSettings(t *testing.T) {
	locale := getLocale(version)
	if locale != `en` {
		t.Fatalf(`expected 'en' but got '%v'`, locale)
	}
}
