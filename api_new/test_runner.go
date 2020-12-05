package api_new

type FileTestRunner struct {
	io          *testIo
	environment *testEnvironment
}

func (r *FileTestRunner) SetUpdateOnly(updateOnly bool) {
	r.environment.updateOnly = updateOnly
}
