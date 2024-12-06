package node

import (
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

type (
	transformValType int

	fakeTransportStruct struct {
		Email    string `faker:"email"`
		Name     string `faker:"name"`
		UnixTime int64  `faker:"unix_time"`
	}

	testTransform        struct{}
	testTransformCfgImpl struct{}
)

const (
	_ transformValType = iota
	transformNil
	transformError
	transformString
	transformStruct
)

func (t *testTransformCfgImpl) GetName() string {
	return "test_transform"
}

func (t *testTransformCfgImpl) GetSourceName() string {
	return "test_transform"
}

func (t *testTransformCfgImpl) GetThreads() int {
	return 1
}

func (t *testTransformCfgImpl) GetOutputBufferSize() int {
	return 1
}

func (t *testTransformCfgImpl) GetInputs() []string {
	return []string{
		"dummy",
	}
}

func (t *testTransform) Init(_ transform.Transform) error {
	return nil
}

func (t *testTransform) Transform(val any) any {
	vType, ok := val.(transformValType)
	if !ok {
		panic("unsupported message type")
	}
	switch vType {
	case transformNil:
		return nil
	case transformError:
		return errors.New(faker.Word())
	case transformString:
		return faker.Word()
	case transformStruct:
		r := fakeTransportStruct{}
		if err := faker.FakeData(&r); err != nil {
			panic(err)
		}
		return &r
	}
	panic("unsupported message type")
}

func TestTransform_Run(t *testing.T) {
	source.RegisterTransform("test_transform", &testTransform{})
	log := zerolog.New(os.Stdout)

	t.Run("multiple input single out", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		transformNode, err := NewTransform(&testTransformCfgImpl{}, &log)
		assert.NoError(t, err)

		inputChan1 := make(chan any, 1)
		transformNode.AddInputChannel(inputChan1)
		inputChan2 := make(chan any, 1)
		transformNode.AddInputChannel(inputChan2)

		outChan := make(chan any, 1)
		transformNode.AddOutputChannel(outChan)

		assert.NoError(t, transformNode.Run())

		countMessages1 := 10
		countMessages2 := 10
		go func() {
			defer close(inputChan1)
			for i := 0; i < countMessages1; i++ {
				inputChan1 <- transformString
			}
		}()

		go func() {
			defer close(inputChan2)
			for i := 0; i < countMessages2; i++ {
				inputChan2 <- transformStruct
			}
		}()

		res := make([]any, 0, countMessages1+countMessages2)
		for val := range outChan {
			res = append(res, val)
		}
		assert.Len(t, res, countMessages1+countMessages2)
	})

	t.Run("single input multiple out", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		transformNode, err := NewTransform(&testTransformCfgImpl{}, &log)
		assert.NoError(t, err)

		inputChan := make(chan any, 1)
		transformNode.AddInputChannel(inputChan)

		outChan1 := make(chan any, 1)
		transformNode.AddOutputChannel(outChan1)
		outChan2 := make(chan any, 1)
		transformNode.AddOutputChannel(outChan2)

		assert.NoError(t, transformNode.Run())

		messagesCount := 10
		go func() {
			defer close(inputChan)
			for i := 0; i < messagesCount; i++ {
				inputChan <- transformString
			}
		}()

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()

			res := make([]any, 0, messagesCount)
			for v := range outChan1 {
				res = append(res, v)
			}
			assert.Len(t, res, messagesCount)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			res := make([]any, 0, messagesCount)
			for v := range outChan2 {
				res = append(res, v)
			}
			assert.Len(t, res, messagesCount)
		}()
		wg.Wait()
	})
}
