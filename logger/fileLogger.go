package logger

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"sync"
)

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
	ReadEvents() (<-chan Event, <-chan error)
	Run()
	Close() error
}

type EventType byte

const (
	_ = iota
	EventDelete
	EventPut
)

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}

type FileTransactionLogger struct {
	events       chan<- Event
	errors       <-chan error
	lastSequence uint64
	file         *os.File
	wg           *sync.WaitGroup
}

func (l *FileTransactionLogger) WriteDelete(key string) {
	l.wg.Add(1)
	l.events <- Event{
		EventType: EventDelete,
		Key:       key,
	}
}

func (l *FileTransactionLogger) WritePut(key, value string) {
	l.wg.Add(1)
	l.events <- Event{
		EventType: EventPut,
		Key:       key,
		Value:     url.QueryEscape(value),
	}
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}
	return &FileTransactionLogger{file: file, wg: &sync.WaitGroup{}}, nil
}

func (l *FileTransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {

		for event := range events {
			l.lastSequence++
			_, err := fmt.Fprintf(l.file, "%d\t%d\t%s\t%s\n",
				l.lastSequence, event.EventType, event.Key, event.Value)
			if err != nil {
				errors <- err
				return
			}
			l.wg.Done()
		}
	}()
}

func (l *FileTransactionLogger) Close() error {
	l.wg.Wait()

	if l.events != nil {
		close(l.events)
	}
	return l.file.Close()
}

func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var event Event
		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			_, _ = fmt.Sscanf(line, "%d\t%d\t%s\t%s",
				&event.Sequence, &event.EventType, &event.Key, &event.Value)

			if l.lastSequence >= event.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			uv, err := url.QueryUnescape(event.Value)
			if err != nil {
				outError <- fmt.Errorf("value decoding failure: %w", err)
				return
			}
			event.Value = uv
			l.lastSequence = event.Sequence
			outEvent <- event
		}
		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()
	return outEvent, outError
}
