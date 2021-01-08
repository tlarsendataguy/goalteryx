package api_new

import (
	"unsafe"
)

type FileTestRunner struct {
	io           *testIo
	plugin       *goPluginSharedMemory
	ayxInterface unsafe.Pointer
	inputs       map[string]string
}

func (r *FileTestRunner) SimulateLifecycle() {
	if len(r.inputs) == 0 {
		simulateInputLifecycle(r.ayxInterface)
	} else {
		return
	}
}

func (r *FileTestRunner) CaptureOutgoingAnchor(name string) *RecordCollector {
	collector := &RecordCollector{}
	sharedMemory := registerTestHarness(collector)

	ii := generateIncomingConnectionInterface()
	callPiAddIncomingConnection(sharedMemory, name, ii)
	callPiAddOutgoingConnection(r.plugin, name, ii)

	return collector
}

func (r *FileTestRunner) ConnectInput(name string, dataFile string) {
	r.inputs[name] = dataFile
}

type RecordCollector struct {
	Config       IncomingRecordInfo
	Name         string
	Data         map[string][]interface{}
	Input        InputConnection
	boolFields   map[string]BoolGetter
	intFields    map[string]IntGetter
	floatFields  map[string]FloatGetter
	stringFields map[string]StringGetter
	timeFields   map[string]TimeGetter
	blobFields   map[string]BytesGetter
}

func (r *RecordCollector) Init(_ Provider) {
	r.Data = make(map[string][]interface{})
	r.boolFields = make(map[string]BoolGetter)
	r.intFields = make(map[string]IntGetter)
	r.floatFields = make(map[string]FloatGetter)
	r.stringFields = make(map[string]StringGetter)
	r.timeFields = make(map[string]TimeGetter)
	r.blobFields = make(map[string]BytesGetter)
}

func (r *RecordCollector) OnInputConnectionOpened(connection InputConnection) {
	r.Input = connection
	r.Name = connection.Name()
	r.Config = connection.Metadata()
	for _, field := range r.Config.Fields() {
		r.Data[field.Name] = []interface{}{}
		switch field.Type {
		case `Bool`:
			boolField, _ := r.Config.GetBoolField(field.Name)
			r.boolFields[field.Name] = boolField.GetValue
		case `Byte`, `Int16`, `Int32`, `Int64`:
			intField, _ := r.Config.GetIntField(field.Name)
			r.intFields[field.Name] = intField.GetValue
		case `Float`, `Double`, `FixedDecimal`:
			floatField, _ := r.Config.GetFloatField(field.Name)
			r.floatFields[field.Name] = floatField.GetValue
		case `String`, `WString`, `V_String`, `V_WString`:
			stringField, _ := r.Config.GetStringField(field.Name)
			r.stringFields[field.Name] = stringField.GetValue
		case `Date`, `DateTime`:
			timeField, _ := r.Config.GetTimeField(field.Name)
			r.timeFields[field.Name] = timeField.GetValue
		case `Blob`, `SpatialObj`:
			blobField, _ := r.Config.GetBlobField(field.Name)
			r.blobFields[field.Name] = blobField.GetValue
		}
	}
}

func (r *RecordCollector) OnRecordPacket(connection InputConnection) {
	packet := connection.Read()
	for packet.Next() {
		record := packet.Record()
		for name, getter := range r.blobFields {
			value := getter(record)
			r.appendDataToField(name, value, value == nil)
		}
		for name, getter := range r.boolFields {
			value, isNull := getter(record)
			r.appendDataToField(name, value, isNull)
		}
		for name, getter := range r.intFields {
			value, isNull := getter(record)
			r.appendDataToField(name, value, isNull)
		}
		for name, getter := range r.floatFields {
			value, isNull := getter(record)
			r.appendDataToField(name, value, isNull)
		}
		for name, getter := range r.stringFields {
			value, isNull := getter(record)
			r.appendDataToField(name, value, isNull)
		}
		for name, getter := range r.timeFields {
			value, isNull := getter(record)
			r.appendDataToField(name, value, isNull)
		}
	}
}

func (r *RecordCollector) OnComplete() {}

func (r *RecordCollector) appendDataToField(fieldName string, value interface{}, isNull bool) {
	if isNull {
		r.Data[fieldName] = append(r.Data[fieldName], nil)
	} else {
		r.Data[fieldName] = append(r.Data[fieldName], value)
	}
}
