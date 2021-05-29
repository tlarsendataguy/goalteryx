package sdk

type Io interface {
	Error(string)
	Warn(string)
	Info(string)
	UpdateProgress(float64)
	DecryptPassword(string) string
}
