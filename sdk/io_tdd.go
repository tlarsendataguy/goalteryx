package sdk

import (
	"fmt"
	"math/rand"
	"time"
)

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

func (t *testIo) UpdateProgress(progress float64) bool {
	println(fmt.Sprintf(`Progress: %v`, progress))
	return true
}

func (t *testIo) DecryptPassword(value string) string {
	return value
}

func (t *testIo) CreateTempFile(ext string) string {
	now := time.Now().Format(`20060102150405`)
	source := rand.NewSource(time.Now().Unix())
	generator := rand.New(source)
	randNum := generator.Intn(1000)
	return fmt.Sprintf(`%v-%03d.%v`, now, randNum, ext)
}

func (t *testIo) NotifyFileInput(message string) {
	println(fmt.Sprintf(`FILE INPUT: %v`, message))
}

func (t *testIo) NotifyFileOutput(message string) {
	println(fmt.Sprintf(`FILE OUTPUT: %v`, message))
}
