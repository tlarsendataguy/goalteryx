package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"goalteryx/recordinfo"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

var recordInfoXml = `<MetaInfo connection="Output">
	<RecordInfo>
		<Field name="ByteField" source="TextInput:" type="Byte"/>
		<Field name="BoolField" source="Formula: 1" type="Bool"/>
		<Field name="Int16Field" source="Formula: 16" type="Int16"/>
		<Field name="Int32Field" source="Formula: 32" type="Int32"/>
		<Field name="Int64Field" source="Formula: 64" type="Int64"/>
		<Field name="FixedDecimalField" scale="6" size="19" source="Formula: 123.45" type="FixedDecimal"/>
		<Field name="FloatField" source="Formula: 678.9" type="Float"/>
		<Field name="DoubleField" source="Formula: 0.12345" type="Double"/>
		<Field name="StringField" size="64" source="Formula: &quot;A&quot;" type="String"/>
		<Field name="WStringField" size="64" source="Formula: &quot;AB&quot;" type="WString"/>
		<Field name="V_StringShortField" size="1000" source="Formula: &quot;ABC&quot;" type="V_String"/>
		<Field name="V_StringLongField" size="2147483647" source="Formula: PadLeft(&quot;&quot;, 500, &apos;B&apos;)" type="V_String"/>
		<Field name="V_WStringShortField" size="10" source="Formula: &quot;XZY&quot;" type="V_WString"/>
		<Field name="V_WStringLongField" size="1073741823" source="Formula: PadLeft(&quot;&quot;, 500, &apos;W&apos;)" type="V_WString"/>
		<Field name="DateField" source="Formula: &apos;2020-01-01&apos;" type="Date"/>
		<Field name="DateTimeField" source="Formula: &apos;2020-02-03 04:05:06&apos;" type="DateTime"/>
	</RecordInfo>
</MetaInfo>
`

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	sourceInfo, _ := recordinfo.FromXml(recordInfoXml)
	_ = sourceInfo.SetIntField(`ByteField`, 1)
	_ = sourceInfo.SetBoolField(`BoolField`, true)
	_ = sourceInfo.SetIntField(`Int16Field`, 1)
	_ = sourceInfo.SetIntField(`Int32Field`, 32)
	_ = sourceInfo.SetIntField(`Int64Field`, 64)
	_ = sourceInfo.SetFloatField(`FixedDecimalField`, 123.45)
	_ = sourceInfo.SetFloatField(`FloatField`, 678.9)
	_ = sourceInfo.SetFloatField(`DoubleField`, 0.12345)
	_ = sourceInfo.SetStringField(`StringField`, `A`)
	_ = sourceInfo.SetStringField(`WStringField`, `AB`)
	_ = sourceInfo.SetStringField(`V_StringShortField`, `ABC`)
	_ = sourceInfo.SetStringField(`V_StringLongField`, strings.Repeat(`B`, 500))
	_ = sourceInfo.SetStringField(`V_WStringShortField`, `XYZ`)
	_ = sourceInfo.SetStringField(`V_WStringLongField`, strings.Repeat(`W`, 500))
	_ = sourceInfo.SetDateField(`DateField`, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))
	_ = sourceInfo.SetDateField(`DateTimeField`, time.Date(2020, 2, 3, 4, 5, 6, 0, time.UTC))

	record, _ := sourceInfo.GenerateRecord()

	info, _ := recordinfo.FromXml(recordInfoXml)

	fieldCount := info.NumFields()
	totalResults := make(map[string]int64, fieldCount)
	doSet := true
	doGenerateRecord := true
	for index := 0; index < fieldCount; index++ {
		field, _ := info.GetFieldByIndex(index)
		start := time.Now()
		for i := 0; i < 10000000; i++ {
			value, err := sourceInfo.GetRawBytesFrom(field.Name, record)
			if err != nil {
				println(fmt.Sprintf(`expected no error getting raw bytes from field %v, but got: %v`, field.Name, err.Error()))
			}
			if doSet {
				err = info.SetFromRawBytes(field.Name, value)
				if err != nil {
					println(fmt.Sprintf(`expected no error setting raw bytes to field %v, but got: %v`, field.Name, err.Error()))
				}
			}
			if doGenerateRecord {
				_, err = info.GenerateRecord()
				if err != nil {
					println(fmt.Sprintf(`expected no error setting raw bytes to field %v, but got: %v`, field.Name, err.Error()))
				}
			}
		}
		end := time.Now()
		duration := end.Sub(start)
		totalResults[field.Name] = duration.Milliseconds()
	}
	resultJson, _ := json.MarshalIndent(totalResults, "", "  ")
	println(fmt.Sprintf(string(resultJson)))
}
