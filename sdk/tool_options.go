package sdk

type toolOptions struct {
	noCache bool
}

type ToolOptionSetter func(toolOptions) toolOptions

func ToolNoCache() ToolOptionSetter {
	return func(options toolOptions) toolOptions {
		options.noCache = true
		return options
	}
}
