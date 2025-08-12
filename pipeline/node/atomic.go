package node

import "github.com/google/uuid"

type AtomicMessage struct {
	id        string
	data      any
	successCh chan<- AtomicMessage
	errorCh   chan<- AtomicMessage
}

func (m AtomicMessage) Success() chan<- AtomicMessage {
	return m.successCh
}

func (m AtomicMessage) Error() chan<- AtomicMessage {
	return m.errorCh
}

func (m AtomicMessage) Data() any {
	return m.data
}

func (m AtomicMessage) ID() string {
	return m.id
}

func NewAtomicMessage(data any, success chan<- AtomicMessage, errors chan<- AtomicMessage) AtomicMessage {
	return AtomicMessage{id: uuid.New().String(), data: data, successCh: success, errorCh: errors}
}
func (m AtomicMessage) MessageFrom(data any) AtomicMessage {
	return AtomicMessage{id: m.ID(), data: data, successCh: m.successCh, errorCh: m.errorCh}
}
