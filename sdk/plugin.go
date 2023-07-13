package sdk

type Plugin interface {
	Init(Provider)
	OnInputConnectionOpened(InputConnection)
	OnRecordPacket(InputConnection)
	OnComplete(nRecordLimit int64)
}
