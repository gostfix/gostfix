package util

import (
	"fmt"
	"sort"
	"strings"
)

var DICT_TYPE_HT = "internal"

type dict_ht struct {
	*Dict
	table map[string]string
	seq   []string
	off   int
}

func dict_ht_lookup(ht *dict_ht) func(string) (string, error) {
	return func(name string) (string, error) {
		if ht.Flags&DICT_FLAG_FOLD_FIX == DICT_FLAG_FOLD_FIX {
			name = strings.ToLower(name)
		}
		val, found := ht.table[name]
		if !found {
			return val, fmt.Errorf("value not found")
		}
		return val, nil
	}
}

func dict_ht_update(ht *dict_ht) func(string, string) int {
	return func(name string, value string) int {
		if ht.Flags&DICT_FLAG_FOLD_FIX == DICT_FLAG_FOLD_FIX {
			name = strings.ToLower(name)
		}
		ht.table[name] = value
		return 0
	}
}

func dict_ht_delete(ht *dict_ht) func(string) int {
	return func(name string) int {
		if ht.Flags&DICT_FLAG_FOLD_FIX == DICT_FLAG_FOLD_FIX {
			name = strings.ToLower(name)
		}
		// we don't get an error if we try to delete something
		// that doesn't exist
		delete(ht.table, name)
		return 0
	}
}

func dict_ht_sequence(ht *dict_ht) func(int, *string, *string) {
	return func(how int, name *string, value *string) {
		switch how {
		case DICT_SEQ_FUN_FIRST:
			// reset the position of our sequence iterator
			ht.off = 0
			// only remake the sequence array if it differs in size from
			// the number of entries in the map
			if len(ht.seq) != len(ht.table) {
				ht.seq = make([]string, len(ht.table))
			}
			// copy all the keys into the sequence array
			i := 0
			for k := range ht.table {
				ht.seq[i] = k
				i++
			}
			// TODO(alf): sort the keys so we have a stable iterator.  I don't
			// know if we need this, it does add some overhead, but gives predictable
			// behaviour
			sort.Strings(ht.seq)
			// fallthrough since getting the next item is exactly the same as getting
			// the first item.
			fallthrough
		case DICT_SEQ_FUN_NEXT:
			// make sure our offset isn't past our array before fetching the next
			// key/val pair.
			if ht.off < len(ht.seq) {
				// get the key out of the sequence array
				*name = ht.seq[ht.off]
				// move the offset to the next position
				ht.off++
				// get the value out of the map based on the key
				// we just pulled out of the sequence array.
				*value = ht.table[*name]
				return
			}
			// If we get here, we are at the end of the sequence list
			*name = ""
			*value = ""
			fallthrough
		default:
			// we are done the iteration, cleanup.
			ht.off = 0
			ht.seq = ht.seq[:0]
		}
	}
}

func dict_ht_close(ht *dict_ht) func() {
	return func() {
		// TODO(alf): I'm not sure if it's worth doing antyhing here
		// we could purge the sequence and map?
	}
}

func DictHtOpen(name string, _ int, flags int) *Dict {
	var ht = dict_ht{
		Dict:  DictAlloc(DICT_TYPE_HT, name),
		table: make(map[string]string),
	}
	ht.Lookup = dict_ht_lookup(&ht)
	ht.Update = dict_ht_update(&ht)
	ht.Delete = dict_ht_delete(&ht)
	ht.Sequence = dict_ht_sequence(&ht)
	ht.Close = dict_ht_close(&ht)

	return ht.Dict
}
