/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package node

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func linesCount(t *testing.T, filePath string) int {
	t.Helper()
	file, err := os.Open(filePath)
	assert.NoError(t, err)
	fileScanner := bufio.NewScanner(file)

	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	return lineCount
}

// fileSinkCollector - write value into separated line in file.
type fileSinkCollector struct {
	m        sync.Mutex
	filePath string
}

func (s *fileSinkCollector) Init(cfg sink.Sink) error {
	fileCfg, ok := cfg.(fileSinkCfg)
	if !ok {
		return errors.New("invalid config")
	}
	s.filePath = fileCfg.GetOutFilePath()
	return nil
}

func (s *fileSinkCollector) Sink(value any) error {
	s.m.Lock()
	defer s.m.Unlock()

	sValue, ok := value.(string)
	if !ok {
		panic(errors.New("invalid sink value value"))
	}
	f, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("open out file: %w", err))
	}
	defer f.Close()
	// write every string on separated lines
	if _, err = f.WriteString(sValue + "\n"); err != nil {
		panic(fmt.Errorf("write out file: %w", err))
	}
	return nil
}

// configuration for fileSinkCollector.
type fileSinkCfgImpl struct {
	outFilePath string
	threads     uint
}

func (s *fileSinkCfgImpl) GetName() string {
	return "test_sink_file"
}

func (s *fileSinkCfgImpl) GetSourceName() string {
	return "test_file"
}

func (s *fileSinkCfgImpl) GetThreads() uint {
	return s.threads
}

func (s *fileSinkCfgImpl) GetInputs() []string {
	return []string{
		"dummy",
	}
}

func (s *fileSinkCfgImpl) GetOutFilePath() string {
	return s.outFilePath
}

type fileSinkCfg interface {
	sink.Sink
	GetOutFilePath() string
}

func TestSink_Run(t *testing.T) {
	source.RegisterSink("test_file", &fileSinkCollector{})
	log := zerolog.New(os.Stdout)
	t.Run("single input", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		outFilePath := filepath.Join(t.TempDir(), time.Now().Format(time.RFC3339Nano))
		fileSink, err := NewSink(
			&fileSinkCfgImpl{
				threads:     1,
				outFilePath: outFilePath,
			},
			&log,
		)
		assert.NoError(t, err)

		writeValues := []string{
			"Pipelaner",
			"Node",
			"Sink",
			"Test",
		}
		inputCh := make(chan any, len(writeValues))
		fileSink.AddInputChannel(inputCh)

		assert.NoError(t, fileSink.Run())

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(inputCh)

			for _, val := range writeValues {
				inputCh <- val
			}
			// wait until messages will be processed by sink
			time.Sleep(3 * time.Second)
		}()

		wg.Wait()
		matches, err := filepath.Glob(outFilePath)
		assert.NoError(t, err)
		assert.True(t, len(matches) == 1)
		assert.Equal(t, len(writeValues), linesCount(t, outFilePath))
	})

	t.Run("multiple inputs", func(t *testing.T) {
		defer goleak.VerifyNone(t)

		outFilePath := filepath.Join(t.TempDir(), time.Now().Format(time.RFC3339Nano))
		fileSink, err := NewSink(
			&fileSinkCfgImpl{
				threads:     1,
				outFilePath: outFilePath,
			},
			&log,
		)
		assert.NoError(t, err)

		writeValues1 := []string{
			"input_1",
			"test_1",
			"values_1",
		}
		inputCh1 := make(chan any, len(writeValues1))
		fileSink.AddInputChannel(inputCh1)

		writeValues2 := []string{
			"input_2",
			"test_2",
			"values_2",
			"!",
		}
		inputCh2 := make(chan any, len(writeValues1))
		fileSink.AddInputChannel(inputCh2)

		assert.NoError(t, fileSink.Run())

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(inputCh1)

			for _, val := range writeValues1 {
				inputCh1 <- val
			}
			// wait until messages will be processed by sink
			time.Sleep(3 * time.Second)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(inputCh2)

			for _, val := range writeValues2 {
				inputCh2 <- val
			}
			// wait until messages will be processed by sink
			time.Sleep(3 * time.Second)
		}()

		wg.Wait()
		matches, err := filepath.Glob(outFilePath)
		assert.NoError(t, err)
		assert.True(t, len(matches) == 1)
		assert.Equal(t, len(writeValues1)+len(writeValues2), linesCount(t, outFilePath))
	})
}
