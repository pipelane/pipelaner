package node

import "github.com/google/uuid"

type AtomicData interface {
	Success() chan<- AtomicData
	Error() chan<- AtomicData
	Data() any
	ID() string
	UpdateData(data any) AtomicData
}

type AtomicMessage struct {
	id        string
	data      any
	successCh chan<- AtomicData
	errorCh   chan<- AtomicData
}

func (m AtomicMessage) Success() chan<- AtomicData {
	return m.successCh
}

func (m AtomicMessage) Error() chan<- AtomicData {
	return m.errorCh
}

func (m AtomicMessage) Data() any {
	return m.data
}

func (m AtomicMessage) ID() string {
	return m.id
}

func NewAtomicMessage(data any, success chan<- AtomicData, errors chan<- AtomicData) AtomicMessage {
	return AtomicMessage{id: uuid.New().String(), data: data, successCh: success, errorCh: errors}
}
func (m AtomicMessage) UpdateData(data any) AtomicData {
	return AtomicMessage{id: m.ID(), data: data, successCh: m.successCh, errorCh: m.errorCh}
}
