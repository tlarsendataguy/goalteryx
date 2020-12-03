package api_new

import "fmt"

type testIo struct{}

func (t *testIo) Error(s string) {
	println(fmt.Sprintf(`ERROR: %v`, s))
}

func (t *testIo) Warn(s string) {
	println(fmt.Sprintf(`WARNING: %v`, s))
}

func (t *testIo) Info(s string) {
	println(fmt.Sprintf(`INFO: %v`, s))
}

func (t *testIo) UpdateProgress(f float64) {
	println(fmt.Sprintf(`Progress: %v`, f))
}
