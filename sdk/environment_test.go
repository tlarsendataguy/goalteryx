package sdk

import (
	"runtime"
	"testing"
)

// These unit tests require Alteryx to be installed.  Change the version below with your version of ALteryx.

const version = `2021.3`

func TestGetSettings(t *testing.T) {
	if runtime.GOOS != `windows` {
		t.Skipf(`TestGetSettings will only pass on Windows with Alteryx installed`)
	}

	locale := getLocale(version)
	if locale != `en` {
		t.Fatalf(`expected 'en' but got '%v'`, locale)
	}
}
