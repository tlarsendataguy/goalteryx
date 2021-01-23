package sdk

import (
	"bufio"
	"fmt"
	"github.com/tlarsen7572/goalteryx/sdk/import_file"
	"os"
	"time"
)

type FileTestRunner struct {
	io     *testIo
	plugin *goPluginSharedMemory
	inputs map[string]*FilePusher
}

func (r *FileTestRunner) SimulateLifecycle() {
	if len(r.inputs) == 0 {
		simulateInputLifecycle(r.plugin.ayxInterface)
	} else {
		for _, pusher := range r.inputs {
			simulateInputLifecycle(pusher.sharedMemory.ayxInterface)
		}
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
	pusher := &FilePusher{file: dataFile}
	sharedMemory := registerTestHarness(pusher)
	pusher.sharedMemory = sharedMemory

	ii := generateIncomingConnectionInterface()
	callPiAddIncomingConnection(r.plugin, name, ii)
	callPiAddOutgoingConnection(sharedMemory, `Output`, ii)

	r.inputs[name] = pusher
}

type FilePusher struct {
	file         string
	sharedMemory *goPluginSharedMemory
	output       OutputAnchor
	provider     Provider
}

func (f *FilePusher) Init(provider Provider) {
	f.provider = provider
	f.output = provider.GetOutputAnchor(`Output`)
}

func (f *FilePusher) OnInputConnectionOpened(_ InputConnection) {
	panic("this should never be called")
}

func (f *FilePusher) OnRecordPacket(_ InputConnection) {
	panic("this should never be called")
}

func (f *FilePusher) OnComplete() {
	file, err := os.Open(f.file)
	if err != nil {
		panic(fmt.Sprintf(`error opening data file: %v`, err.Error()))
	}

	scanner := bufio.NewScanner(file)
	success := scanner.Scan()
	if !success {
		return
	}
	fieldNames := import_file.Preprocess(scanner.Bytes())
	success = scanner.Scan()
	if !success {
		return
	}
	fieldTypes := import_file.Preprocess(scanner.Bytes())

	extractor := import_file.NewExtractor(fieldNames, fieldTypes)
	infoEditor := &EditingRecordInfo{}
	source := `FilePusher`

	for _, field := range extractor.Fields() {
		switch field.Type {
		case `Bool`:
			infoEditor.AddBoolField(field.Name, source)
		case `Byte`:
			infoEditor.AddByteField(field.Name, source)
		case `Int16`:
			infoEditor.AddInt16Field(field.Name, source)
		case `Int32`:
			infoEditor.AddInt32Field(field.Name, source)
		case `Int64`:
			infoEditor.AddInt64Field(field.Name, source)
		case `Float`:
			infoEditor.AddFloatField(field.Name, source)
		case `Double`:
			infoEditor.AddDoubleField(field.Name, source)
		case `FixedDecimal`:
			infoEditor.AddFixedDecimalField(field.Name, source, field.Size, field.Scale)
		case `Date`:
			infoEditor.AddDateField(field.Name, source)
		case `DateTime`:
			infoEditor.AddDateTimeField(field.Name, source)
		case `String`:
			infoEditor.AddStringField(field.Name, source, field.Size)
		case `WString`:
			infoEditor.AddWStringField(field.Name, source, field.Size)
		case `V_String`:
			infoEditor.AddV_StringField(field.Name, source, field.Size)
		case `V_WString`:
			infoEditor.AddV_WStringField(field.Name, source, field.Size)
		case `Blob`:
			infoEditor.AddBlobField(field.Name, source, field.Size)
		case `SpatialObj`:
			infoEditor.AddSpatialObjField(field.Name, source, field.Size)
		}
	}
	outInfo := infoEditor.GenerateOutgoingRecordInfo()
	f.output.Open(outInfo)

	for scanner.Scan() {
		preprocessed := import_file.Preprocess(scanner.Bytes())
		data := extractor.Extract(preprocessed)
		for fieldName, value := range data.BlobFields {
			if value == nil {
				outInfo.BlobFields[fieldName].SetNullBlob()
			} else {
				outInfo.BlobFields[fieldName].SetBlob(value.([]byte))
			}
		}
		for fieldName, value := range data.BoolFields {
			if value == nil {
				outInfo.BoolFields[fieldName].SetNullBool()
			} else {
				outInfo.BoolFields[fieldName].SetBool(value.(bool))
			}
		}
		for fieldName, value := range data.IntFields {
			if value == nil {
				outInfo.IntFields[fieldName].SetNullInt()
			} else {
				outInfo.IntFields[fieldName].SetInt(value.(int))
			}
		}
		for fieldName, value := range data.DecimalFields {
			if value == nil {
				outInfo.FloatFields[fieldName].SetNullFloat()
			} else {
				outInfo.FloatFields[fieldName].SetFloat(value.(float64))
			}
		}
		for fieldName, value := range data.DateTimeFields {
			if value == nil {
				outInfo.DateTimeFields[fieldName].SetNullDateTime()
			} else {
				outInfo.DateTimeFields[fieldName].SetDateTime(value.(time.Time))
			}
		}
		for fieldName, value := range data.StringFields {
			if value == nil {
				outInfo.StringFields[fieldName].SetNullString()
			} else {
				outInfo.StringFields[fieldName].SetString(value.(string))
			}
		}
		f.output.Write()
	}
	f.output.UpdateProgress(1.0)
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
			if value == nil {
				r.appendDataToField(name, value, true)
			} else {
				copyValue := make([]byte, len(value))
				copy(copyValue, value)
				r.appendDataToField(name, copyValue, false)
			}
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
