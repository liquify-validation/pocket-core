package types

import (
	"bytes"
	"github.com/tendermint/tendermint/libs/kv"
)

// Iterator over all the keys with a certain prefix in ascending order
func KVStorePrefixIterator(kvs KVStore, prefix []byte) (Iterator, error) {
	return kvs.Iterator(prefix, PrefixEndBytes(prefix))
}

// Iterator over all the keys with a certain prefix in descending order.
func KVStoreReversePrefixIterator(kvs KVStore, prefix []byte) (Iterator, error) {
	return kvs.ReverseIterator(prefix, PrefixEndBytes(prefix))
}

// Compare two KVstores, return either the first key/value pair
// at which they differ and whether or not they are equal, skipping
// value comparison for a set of provided prefixes
func DiffKVStores(a KVStore, b KVStore, prefixesToSkip [][]byte) (kvA kv.Pair, kvB kv.Pair, count int64, equal bool) {
	iterA, _ := a.Iterator(nil, nil)
	iterB, _ := b.Iterator(nil, nil)
	count = int64(0)
	for {
		if !iterA.Valid() && !iterB.Valid() {
			break
		}
		var kvA, kvB kv.Pair
		if iterA.Valid() {
			kvA = kv.Pair{Key: iterA.Key(), Value: iterA.Value()}
			iterA.Next()
		}
		if iterB.Valid() {
			kvB = kv.Pair{Key: iterB.Key(), Value: iterB.Value()}
			iterB.Next()
		}
		if !bytes.Equal(kvA.Key, kvB.Key) {
			return kvA, kvB, count, false
		}
		compareValue := true
		for _, prefix := range prefixesToSkip {
			// Skip value comparison if we matched a prefix
			if bytes.Equal(kvA.Key[:len(prefix)], prefix) {
				compareValue = false
			}
		}
		if compareValue && !bytes.Equal(kvA.Value, kvB.Value) {
			return kvA, kvB, count, false
		}
		count++
	}
	return kv.Pair{}, kv.Pair{}, count, true
}

// PrefixEndBytes returns the []byte that would end a
// range query for all []byte with a certain prefix
// Deals with last byte of prefix being FF without overflowing
func PrefixEndBytes(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}

	end := make([]byte, len(prefix))
	copy(end, prefix)

	for {
		if end[len(end)-1] != byte(255) {
			end[len(end)-1]++
			break
		} else {
			end = end[:len(end)-1]
			if len(end) == 0 {
				end = nil
				break
			}
		}
	}
	return end
}

// InclusiveEndBytes returns the []byte that would end a
// range query such that the input would be included
func InclusiveEndBytes(inclusiveBytes []byte) (exclusiveBytes []byte) {
	exclusiveBytes = append(inclusiveBytes, byte(0x00))
	return exclusiveBytes
}

//----------------------------------------
func Cp(bz []byte) (ret []byte) {
	if bz == nil {
		return nil
	}
	ret = make([]byte, len(bz))
	copy(ret, bz)
	return ret
}