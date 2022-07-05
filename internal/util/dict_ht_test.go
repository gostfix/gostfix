package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDictHt(t *testing.T) {
	var dict *Dict = DictHtOpen("test_table", 0, DICT_FLAG_NONE)

	dict.Update("key1", "val1")
	val, err := dict.Lookup("key1")
	assert.NoError(t, err)
	assert.Equal(t, "val1", val)
}

func TestDictHtClosure(t *testing.T) {
	var dict01 *Dict = DictHtOpen("test_table01", 0, DICT_FLAG_NONE)
	var dict02 *Dict = DictHtOpen("test_table02", 0, DICT_FLAG_NONE)
	var dict03 *Dict = DictHtOpen("test_table03", 0, DICT_FLAG_NONE)
	var dict04 *Dict = DictHtOpen("test_table04", 0, DICT_FLAG_NONE)

	dict01.Update("key1", "val1")
	dict02.Update("key2", "val2")
	dict03.Update("key3", "val3")
	dict04.Update("key4", "val4")

	val, err := dict01.Lookup("key1")
	assert.NoError(t, err)
	assert.Equal(t, "val1", val)
	val, err = dict02.Lookup("key2")
	assert.NoError(t, err)
	assert.Equal(t, "val2", val)
	val, err = dict03.Lookup("key3")
	assert.NoError(t, err)
	assert.Equal(t, "val3", val)
	val, err = dict04.Lookup("key4")
	assert.NoError(t, err)
	assert.Equal(t, "val4", val)

}

func TestDictHtSequence(t *testing.T) {
	var dict *Dict = DictHtOpen("test_table", 0, DICT_FLAG_NONE)

	dict.Update("key01", "val01")
	dict.Update("key02", "val02")
	dict.Update("key03", "val03")
	dict.Update("key04", "val04")
	dict.Update("key05", "val05")
	dict.Update("key06", "val06")
	dict.Update("key07", "val07")
	dict.Update("key08", "val08")
	dict.Update("key09", "val09")
	dict.Update("key10", "val10")
	dict.Update("key11", "val11")

	var key1, val1 string
	dict.Sequence(DICT_SEQ_FUN_FIRST, &key1, &val1)

	var key2, val2 string
	dict.Sequence(DICT_SEQ_FUN_FIRST, &key2, &val2)

	assert.Equal(t, key1, key2)
	assert.Equal(t, val1, val2)
}
