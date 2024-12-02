package node

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

type (
	inputValType int

	fakeInputStruct struct {
		Email string `faker:"email"`
		Name  string `faker:"name"`
	}

	testInput struct {
		generateCount int
		generateType  inputValType
	}

	testInputCfg interface {
		input.Input
		GetGenerateCount() int
		GetGenerateType() inputValType
	}
	testInputCfgImpl struct {
		generateCount int
		generateType  inputValType
	}
)

const (
	_ inputValType = iota
	inputNil
	inputError
	inputStruct
)

func (t *testInputCfgImpl) GetName() string {
	return "test_input"
}

func (t *testInputCfgImpl) GetSourceName() string {
	return "test_input"
}

func (t *testInputCfgImpl) GetThreads() int {
	return 1
}

func (t *testInputCfgImpl) GetOutputBufferSize() int {
	return 1
}

func (t *testInputCfgImpl) GetInputs() []string {
	return []string{
		"dummy",
	}
}

func (t *testInputCfgImpl) GetGenerateCount() int {
	return t.generateCount
}

func (t *testInputCfgImpl) GetGenerateType() inputValType {
	return t.generateType
}

func (t *testInput) Init(cfg input.Input) error {
	tCfg, ok := cfg.(testInputCfg)
	if !ok {
		panic("invalid config type")
	}
	t.generateCount = tCfg.GetGenerateCount()
	t.generateType = tCfg.GetGenerateType()
	return nil
}

func (t *testInput) Generate(_ context.Context, input chan<- any) {
	switch t.generateType {
	case inputNil:
		input <- nil
	case inputError:
		input <- fmt.Errorf("dummy")
	case inputStruct:
		for i := 0; i < t.generateCount; i++ {
			r := fakeInputStruct{}
			if err := faker.FakeData(&r); err != nil {
				panic(err)
			}
			input <- r
		}
	}
}

func TestInput_Run(t *testing.T) {
	source.RegisterInput("test_input", &testInput{})
	var logBuffer bytes.Buffer
	log := zerolog.New(&logBuffer)

	t.Run("generate nil message", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		defer logBuffer.Reset()

		inputNode, err := NewInput(&testInputCfgImpl{
			generateType: inputNil,
		}, &log)
		assert.NoError(t, err)

		outChan := make(chan any, 1)
		inputNode.AddOutputChannel(outChan)

		assert.NoError(t, inputNode.Run(context.Background()))
		<-outChan
		assert.True(t, bytes.Contains(logBuffer.Bytes(), []byte("received nil message")))
	})

	t.Run("generate error message", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		defer logBuffer.Reset()

		inputNode, err := NewInput(&testInputCfgImpl{
			generateType: inputError,
		}, &log)
		assert.NoError(t, err)

		outChan := make(chan any, 1)
		inputNode.AddOutputChannel(outChan)

		assert.NoError(t, inputNode.Run(context.Background()))
		<-outChan
		assert.True(t, bytes.Contains(logBuffer.Bytes(), []byte("received error")))
	})

	t.Run("generate message, multiple output", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		defer logBuffer.Reset()

		messagesCount := 3
		inputNode, err := NewInput(&testInputCfgImpl{
			generateType:  inputStruct,
			generateCount: messagesCount,
		}, &log)
		assert.NoError(t, err)

		outChan1 := make(chan any, messagesCount)
		inputNode.AddOutputChannel(outChan1)
		outChan2 := make(chan any, messagesCount)
		inputNode.AddOutputChannel(outChan2)

		assert.NoError(t, inputNode.Run(context.Background()))

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()

			r := make([]any, 0, messagesCount)
			for v := range outChan1 {
				r = append(r, v)
			}
			assert.Len(t, r, messagesCount)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			r := make([]any, 0, messagesCount)
			for v := range outChan2 {
				r = append(r, v)
			}
			assert.Len(t, r, messagesCount)
		}()

		wg.Wait()
	})
}
