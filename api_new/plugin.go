package api_new

type Plugin interface {
	Init(Provider)
	OnInputConnectionOpened(InputConnection)
	OnRecordPacket(InputConnection)
	OnComplete()
}
