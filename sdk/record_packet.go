package sdk

import (
	"unsafe"
)

type RecordPacket interface {
	Next() bool
	Record() Record
}

type RecordCache unsafe.Pointer
type Record = unsafe.Pointer

func NewRecordPacket(cache RecordCache, size int, fixedLen int, hasVarData bool) RecordPacket {
	return &impRecordPacket{
		cache:           cache,
		size:            uintptr(size),
		currentPosition: 0,
		fixedLen:        uintptr(fixedLen),
		hasVarData:      hasVarData,
		currentRecord:   nil,
	}
}

type impRecordPacket struct {
	cache           RecordCache
	size            uintptr
	currentPosition uintptr
	fixedLen        uintptr
	hasVarData      bool
	currentRecord   Record
}

func (p *impRecordPacket) Next() bool {
	if p.size == 0 {
		return false
	}
	if p.atFirstRecord() {
		p.currentRecord = Record(p.cache)
		return true
	}
	p.currentPosition += p.fixedLen
	if p.afterLastRecord() {
		p.currentRecord = nil
		return false
	}
	if p.hasVarData {
		varSize := *(*uint32)(unsafe.Pointer(uintptr(p.cache) + p.currentPosition))
		p.currentPosition += 4 + uintptr(varSize)
	}
	if p.afterLastRecord() {
		p.currentRecord = nil
		return false
	}
	p.currentRecord = Record(uintptr(p.cache) + p.currentPosition)
	return true
}

func (p *impRecordPacket) Record() Record {
	return p.currentRecord
}

func (p *impRecordPacket) atFirstRecord() bool {
	return p.currentRecord == nil && p.currentPosition == 0
}

func (p *impRecordPacket) afterLastRecord() bool {
	return p.currentPosition >= p.size
}
