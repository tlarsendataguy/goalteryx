package api_new

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

type provider struct {
	sharedMemory  *goPluginSharedMemory
	io            Io
	environment   Environment
	outputAnchors map[string]*outputAnchor
}

func (p *provider) ToolConfig() string {
	return utf16PtrToString(p.sharedMemory.toolConfig, int(p.sharedMemory.toolConfigLen))
}

func (p *provider) Io() Io {
	return p.io
}

func (p *provider) GetOutputAnchor(name string) OutputAnchor {
	anchor, ok := p.outputAnchors[name]
	if ok {
		return anchor
	}
	anchorData := getOrCreateOutputAnchor(p.sharedMemory, name)
	anchor = &outputAnchor{data: anchorData}
	p.outputAnchors[name] = anchor
	return anchor
}

func (p *provider) Environment() Environment {
	return p.environment
}

type outputAnchor struct {
	data     *goOutputAnchorData
	metaData *OutgoingRecordInfo
}

func (a *outputAnchor) Name() string {
	name := utf16PtrToString(a.data.name, utf16PtrLen(a.data.name))
	return name
}

func (a *outputAnchor) IsOpen() bool {
	return a.data.isOpen == 1
}

func (a *outputAnchor) Metadata() *OutgoingRecordInfo {
	return a.metaData
}

func (a *outputAnchor) Open(info *OutgoingRecordInfo) {
	a.metaData = info
	a.data.fixedSize = uint32(info.FixedSize())
	if info.HasVarFields() {
		a.data.hasVarFields = 1
	}
	a.data.recordCache = allocateCache(cacheSize)
	a.data.recordCacheSize = uint32(cacheSize)
	xmlStr := info.toXml(a.Name())
	openOutgoingAnchor(a.data, xmlStr)
}

func (a *outputAnchor) Write() {
	nextRecordSize := a.metaData.DataSize()
	dataStartsAt := int(a.data.recordCachePosition)
	currentFixedPosition := 0
	currentVarPosition := int(a.data.fixedSize) + 4
	varLen := 0
	hasVar := false
	if dataStartsAt+nextRecordSize > cacheSize {
		callWriteRecords(unsafe.Pointer(a.data))
		dataStartsAt = 0
	}
	cache := ptrToBytes(a.data.recordCache, dataStartsAt, nextRecordSize)
	for _, field := range a.metaData.outgoingFields {
		if field.isFixedLen {
			copy(cache[currentFixedPosition:], field.CurrentValue)
			currentFixedPosition += len(field.CurrentValue)
			continue
		}
		hasVar = true
		if field.CurrentValue[0] == 1 {
			copy(cache[currentFixedPosition:], []byte{1, 0, 0, 0})
			currentFixedPosition += 4
			continue
		}
		if len(field.CurrentValue) == 0 {
			copy(cache[currentFixedPosition:], []byte{0, 0, 0, 0})
			currentFixedPosition += 4
			continue
		}
		varWritten := varBytesToCache(field.CurrentValue[1:], cache, currentFixedPosition, currentVarPosition)
		currentFixedPosition += 4
		varLen += varWritten
		currentVarPosition += varWritten
	}
	if hasVar {
		binary.LittleEndian.PutUint32(cache[currentFixedPosition:currentFixedPosition+4], uint32(varLen))
		currentFixedPosition += 4
	}
	if varLen+currentFixedPosition != nextRecordSize {
		panic(fmt.Sprintf(`mismatch between actual write of %v and calculated write of %v`, varLen+currentFixedPosition, nextRecordSize))
	}
	a.data.recordCachePosition += uint32(nextRecordSize)
}

func (a *outputAnchor) UpdateProgress(progress float64) {
	panic("implement me")
}

func varBytesToCache(varBytes []byte, cache []byte, fixedPosition int, varPosition int) int {
	varWritten := len(varBytes)
	varDataLen := uint32(varWritten)

	// Small string optimization
	if varDataLen < 4 {
		varDataLen <<= 28
		fixedBytes := make([]byte, 4)
		copy(fixedBytes, varBytes)
		varDataUint32 := binary.LittleEndian.Uint32(fixedBytes) | varDataLen
		binary.LittleEndian.PutUint32(cache[fixedPosition:fixedPosition+4], varDataUint32)
		return 0
	}

	binary.LittleEndian.PutUint32(cache[fixedPosition:fixedPosition+4], uint32(varPosition-fixedPosition))

	if varDataLen < 128 {
		cache[varPosition] = byte(varDataLen*2) | 1 // Alteryx seems to multiply all var lens by 2
		varPosition += 1
		varWritten += 1
	} else {
		binary.LittleEndian.PutUint32(cache[varPosition:varPosition+4], varDataLen*2) // Alteryx seems to multiply all var lens by 2
		varPosition += 4
		varWritten += 4
	}

	copy(cache[varPosition:], varBytes)
	return varWritten
}
