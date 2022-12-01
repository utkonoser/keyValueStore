package main

import (
	"os"
	"testing"
)

func fileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func TestNewFileTransactionLogger(t *testing.T) {
	const filename = "tmp/test-logger.txt"
	defer os.Remove(filename)

	t1, err := NewFileTransactionLogger(filename)

	if t1 == nil {
		t.Error("Logger is nil?")
	}

	if err != nil {
		t.Errorf(err.Error())
	}

	if !fileExist(filename) {
		t.Errorf("File %s doesnot exist!", filename)
	}
}
