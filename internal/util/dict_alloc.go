package util

import (
	"fmt"
	"time"
)

func dict_lookup(dict *Dict) func(string) (string, error) {
	return func(name string) (string, error) {
		MsgFatal("lookup operation not supported", "type", dict.Type, "name", dict.Name)
		return "", fmt.Errorf("lookup operation not supported")
	}
}

func dict_update(dict *Dict) func(string, string) int {
	return func(name string, value string) int {
		MsgFatal("update operation not supported", "type", dict.Type, "name", dict.Name)
		return 0
	}
}

func dict_delete(dict *Dict) func(string) int {
	return func(name string) int {
		MsgFatal("delete operation not supported", "type", dict.Type, "name", dict.Name)
		return 0
	}
}

func dict_sequence(dict *Dict) func(int, *string, *string) {
	return func(how int, name *string, value *string) {
		MsgFatal("sequence operation not supported", "type", dict.Type, "name", dict.Name)
	}
}

func dict_lock(dict *Dict) func(int) int {
	return func(operation int) int {
		if dict.LockFd < 0 {
			return MyFlock(dict.LockFd, dict.LockType, operation)
		}
		return 0
	}
}

func dict_close(dict *Dict) func() {
	return func() {
		MsgFatal("close operation not supported", "type", dict.Type, "name", dict.Name)
	}
}

func DictAlloc(dict_type string, dict_name string) *Dict {
	var dict = Dict{
		Type:     dict_type,
		Name:     dict_name,
		Flags:    DICT_FLAG_FIXED,
		Lookup:   nil,
		Update:   nil,
		Delete:   nil,
		Sequence: nil,
		Lock:     nil,
		Close:    nil,
		LockType: 0,
		LockFd:   -1,
		StatFd:   -1,
		ModTime:  time.Time{},
		Owner:    DictOwner{Status: -1, Uid: -1},
		Error:    DICT_ERR_NONE,
	}

	dict.Lookup = dict_lookup(&dict)
	dict.Update = dict_update(&dict)
	dict.Delete = dict_delete(&dict)
	dict.Sequence = dict_sequence(&dict)
	dict.Lock = dict_lock(&dict)
	dict.Close = dict_close(&dict)
	return &dict
}
