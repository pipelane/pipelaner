package node

import (
	"os"
	"sync"
	"testing"

	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

type testSinkAtomic struct{}

func (s *testSinkAtomic) Init(_ sink.Sink) error {
	return nil
}

func (s *testSinkAtomic) Sink(val any) error {
	switch v := val.(type) {
	case AtomicMessage:
		v.Success() <- v
		return nil
	default:
		panic("unsupported message type")
	}
}

func (s *testSinkAtomic) GetName() string {
	return "test_atomic_sink"
}

func (s *testSinkAtomic) GetSourceName() string {
	return "test_atomic_sink"
}

func (s *testSinkAtomic) GetThreads() uint {
	return 1
}

func (s *testSinkAtomic) GetInputs() []string {
	return []string{
		"dummy",
	}
}

type testTransformAtomic struct{}

func (t *testTransformAtomic) GetName() string {
	return "test_atomic_transform"
}

func (t *testTransformAtomic) GetSourceName() string {
	return "test_atomic_transform"
}

func (t *testTransformAtomic) GetInputs() []string {
	return []string{
		"dummy",
	}
}

func (t *testTransformAtomic) Init(_ transform.Transform) error {
	return nil
}

func (t *testTransformAtomic) Transform(val any) any {
	switch v := val.(type) {
	case AtomicMessage:
		vals := t.Transform(v.Data())
		if vs, ok := vals.(error); ok {
			return vs
		}
		newV := v.MessageFrom(vals)
		return newV
	case int:
		v++
		return v
	default:
		panic("unsupported message type")
	}
}

func (t *testTransformAtomic) GetThreads() uint {
	return 1
}

func (t *testTransformAtomic) GetOutputBufferSize() uint {
	return 1
}

func TestNodes_Run(t *testing.T) {
	source.RegisterTransform("test_atomic_transform", &testTransformAtomic{})
	source.RegisterSink("test_atomic_sink", &testSinkAtomic{})
	log := zerolog.New(os.Stdout)

	t.Run("single input multiple out", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		transformNode, err := NewTransform(&testTransformAtomic{}, &log)
		assert.NoError(t, err)

		sinkNode, err := NewSink(&testSinkAtomic{}, &log)
		assert.NoError(t, err)
		messagesCount := 10

		inputChan := make(chan any, messagesCount)
		transformNode.AddInputChannel(inputChan)

		outChan1 := make(chan any, messagesCount)
		transformNode.AddOutputChannel(outChan1)

		sinkNode.AddInputChannel(outChan1)

		assert.NoError(t, transformNode.Run())
		assert.NoError(t, sinkNode.Run())
		success := make(chan AtomicMessage, messagesCount)
		ids := map[string]int{}
		go func() {
			for i := 0; i < messagesCount; i++ {
				message := NewAtomicMessage(i, success, nil)
				ids[message.ID()] = i + 1
				inputChan <- message
			}
			close(inputChan)
		}()
		res := make([]AtomicMessage, 0, messagesCount)
		counter := 0
		for v := range success {
			res = append(res, v)
			counter++
			if counter == messagesCount {
				close(success)
			}
		}
		for _, v := range res {
			data, ok := v.Data().(int)
			assert.True(t, ok)
			assert.Equal(t, ids[v.ID()], data)
		}
		assert.Len(t, res, messagesCount)
	})

	t.Run("multiple input single out", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		transformNode, err := NewTransform(&testTransformAtomic{}, &log)
		assert.NoError(t, err)

		sinkNode, err := NewSink(&testSinkAtomic{}, &log)
		assert.NoError(t, err)
		messagesCount := 10

		inputChan1 := make(chan any, messagesCount)
		transformNode.AddInputChannel(inputChan1)

		inputChan2 := make(chan any, messagesCount)
		transformNode.AddInputChannel(inputChan2)

		outChan1 := make(chan any, messagesCount)
		transformNode.AddOutputChannel(outChan1)

		sinkNode.AddInputChannel(outChan1)

		assert.NoError(t, transformNode.Run())
		assert.NoError(t, sinkNode.Run())
		success := make(chan AtomicMessage, messagesCount)
		ids := map[string]int{}
		mx := sync.RWMutex{}
		go func() {
			for i := 0; i < messagesCount; i++ {
				message := NewAtomicMessage(i, success, nil)
				mx.Lock()
				ids[message.ID()] = i + 1
				mx.Unlock()
				inputChan1 <- message
			}
			close(inputChan1)
		}()
		go func() {
			for i := 0; i < messagesCount; i++ {
				message := NewAtomicMessage(i, success, nil)
				mx.Lock()
				ids[message.ID()] = i + 1
				mx.Unlock()
				inputChan2 <- message
			}
			close(inputChan2)
		}()
		res := make([]AtomicMessage, 0, messagesCount*2)
		counter := 0
		for v := range success {
			res = append(res, v)
			counter++
			if counter == messagesCount*2 {
				close(success)
			}
		}
		for _, v := range res {
			data, ok := v.Data().(int)
			assert.True(t, ok)
			assert.Equal(t, ids[v.ID()], data)
		}
		assert.Len(t, res, messagesCount*2)
	})
}
