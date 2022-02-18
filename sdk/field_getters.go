package sdk

import (
	"encoding/binary"
	"math"
	"time"
	"unicode/utf16"
)

type IntGetter func(Record) (int, bool)
type FloatGetter func(Record) (float64, bool)
type BoolGetter func(Record) (bool, bool)
type TimeGetter func(Record) (time.Time, bool)
type StringGetter func(Record) (string, bool)

func bytesToByte(getBytes BytesGetter) IntGetter {
	return func(record Record) (int, bool) {
		bytes := getBytes(record)
		if bytes[1] == 1 {
			return 0, true
		}
		return int(bytes[0]), false
	}
}

func bytesToInt16(getBytes BytesGetter) IntGetter {
	return func(record Record) (int, bool) {
		bytes := getBytes(record)
		if bytes[2] == 1 {
			return 0, true
		}
		return int(int16(binary.LittleEndian.Uint16(bytes))), false
	}
}

func bytesToInt32(getBytes BytesGetter) IntGetter {
	return func(record Record) (int, bool) {
		bytes := getBytes(record)
		if bytes[4] == 1 {
			return 0, true
		}
		return int(int32(binary.LittleEndian.Uint32(bytes))), false
	}
}

func bytesToInt64(getBytes BytesGetter) IntGetter {
	return func(record Record) (int, bool) {
		bytes := getBytes(record)
		if bytes[8] == 1 {
			return 0, true
		}
		return int(int64(binary.LittleEndian.Uint64(bytes))), false
	}
}

func bytesToFloat(getBytes BytesGetter) FloatGetter {
	return func(record Record) (float64, bool) {
		bytes := getBytes(record)
		if bytes[4] == 1 {
			return 0, true
		}
		return float64(math.Float32frombits(binary.LittleEndian.Uint32(bytes))), false
	}
}

func bytesToDouble(getBytes BytesGetter) FloatGetter {
	return func(record Record) (float64, bool) {
		bytes := getBytes(record)
		if bytes[8] == 1 {
			return 0, true
		}
		return math.Float64frombits(binary.LittleEndian.Uint64(bytes)), false
	}
}

func bytesToString(getBytes BytesGetter, size int) StringGetter {
	return func(record Record) (string, bool) {
		bytes := getBytes(record)
		if bytes[size] == 1 {
			return ``, true
		}
		return string(truncateAtNullByte(bytes)), false
	}
}

func truncateAtNullUtf16(raw []uint16) []uint16 {
	var dataLen int
	for dataLen = 0; dataLen < len(raw); dataLen++ {
		if raw[dataLen] == 0 {
			break
		}
	}
	return raw[:dataLen]
}

func bytesToWString(getBytes BytesGetter, size int) StringGetter {
	return func(record Record) (string, bool) {
		bytes := getBytes(record)
		if bytes[size*2] == 1 {
			return ``, true
		}
		utf16Bytes := bytesToUtf16(bytes)
		utf16Bytes = truncateAtNullUtf16(utf16Bytes)
		if len(utf16Bytes) == 0 {
			return ``, false
		}
		value := string(utf16.Decode(utf16Bytes))
		return value, false
	}
}

func bytesToV_String(getBytes BytesGetter, _ int) StringGetter {
	return func(record Record) (string, bool) {
		bytes := getBytes(record)
		if bytes == nil {
			return ``, true
		}
		return string(bytes), false
	}
}

func bytesToV_WString(getBytes BytesGetter, _ int) StringGetter {
	return func(record Record) (string, bool) {
		bytes := getBytes(record)
		if bytes == nil {
			return ``, true
		}
		if len(bytes) == 0 {
			return ``, false
		}
		utf16Bytes := bytesToUtf16(bytes)
		value := string(utf16.Decode(utf16Bytes))
		return value, false
	}
}
