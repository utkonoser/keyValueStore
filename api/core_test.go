package api

import (
	"errors"
	"testing"
)

func TestPut(t *testing.T) {
	const key = "test_key"
	const value = "test_value"

	var contains bool
	var val interface{}

	defer delete(store.m, key)

	// sanity check
	_, contains = store.m[key]
	if contains {
		t.Error("key/value already exist")
	}

	// err should be nil
	err := Put(key, value)
	if err != nil {
		t.Error(err)
	}

	val, contains = store.m[key]
	if !contains {
		t.Error("create failed!")
	}
	if val != value {
		t.Error("val/value mismatch!")
	}
}

func TestGet(t *testing.T) {
	const key = "test_key"
	const value = "test_value"

	var err error
	var val interface{}

	defer delete(store.m, key)
	// Read a non-thing
	val, err = Get(key)
	if err == nil {
		t.Error("expected an error!")
	}
	if !errors.Is(err, ErrorNoSuchKey) {
		t.Error("unexpected an error!")
	}

	store.m[key] = value
	// Read a non-empty
	val, err = Get(key)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if val != value {
		t.Error("val/value mismatch")
	}
}

func TestDelete(t *testing.T) {
	const key = "test_key"
	const value = "test_value"

	var contains bool

	defer delete(store.m, key)

	store.m[key] = value
	_, contains = store.m[key]
	if !contains {
		t.Error("key/value doesnt exist!")
	}

	_ = Delete(key)
	_, contains = store.m[key]
	if contains {
		t.Error("delete failed!")
	}

}
