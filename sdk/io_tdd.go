package sdk

import "fmt"

type testIo struct{}

func (t *testIo) Error(message string) {
	println(fmt.Sprintf(`ERROR: %v`, message))
}

func (t *testIo) Warn(message string) {
	println(fmt.Sprintf(`WARNING: %v`, message))
}

func (t *testIo) Info(message string) {
	println(fmt.Sprintf(`INFO: %v`, message))
}

func (t *testIo) UpdateProgress(progress float64) {
	println(fmt.Sprintf(`Progress: %v`, progress))
}
