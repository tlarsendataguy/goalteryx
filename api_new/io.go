package api_new

type Io interface {
	Error(string)
	Warn(string)
	Info(string)
	UpdateProgress(float64)
}
