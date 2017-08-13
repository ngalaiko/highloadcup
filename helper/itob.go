package helper

import (
	"encoding/binary"
)

// Itob returns an 8-byte big endian representation of v
func Itob(v uint32) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
