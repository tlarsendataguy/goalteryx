package api_new

type Environment interface {
	UpdateOnly() bool
	UpdateMode() string
	DesignerVersion() string
	WorkflowDir() string
	AlteryxInstallDir() string
	AlteryxLocale() string
	ToolId() int
	UpdateToolConfig(string)
}
