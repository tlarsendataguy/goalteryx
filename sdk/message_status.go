package sdk

type MessageStatus int

const (
	Info                          MessageStatus = 1
	TransientInfo                 MessageStatus = 0x40000000 | 1
	Warning                       MessageStatus = 2
	TransientWarning              MessageStatus = 0x40000000 | 2
	Error                         MessageStatus = 3
	Complete                      MessageStatus = 4
	FieldConversionError          MessageStatus = 5
	TransientFieldConversionError MessageStatus = 0x40000000 | 5
	FileInput                     MessageStatus = 8
	FileOutput                    MessageStatus = 9
	UpdateOutputMetaInfoXml       MessageStatus = 10
	RecordCountString             MessageStatus = 50
	BrowseEverywhereFileName      MessageStatus = 70
)
