package pickle

import (
	"encoding/binary"
	"errors"
	"fmt"
	"unicode/utf16"

	log "github.com/sirupsen/logrus"
)

//Pickles are not strictly python pickles, but Google-made, and similar
// they map to concrete types, are prefixed with lengths to read
// and are uint32-aligned data (data will always occupy 4-byte blocks)
// so strings may be padded

type PickleType int

const (
	PickleBool PickleType = iota
	PickleInt16
	PickleUint16
	PickleInt32
	PickleUint32
	PickleInt64
	PickleUint64
	PickleFloat //single
	PickleDouble
	PickleBlob
	PickleString
	PickleUTF16
	PickleDateTime
	PicklePickle
	PickleUTF16Count
)

// Reads a Chrome-Pickle string.
// returns the next index of the passed []byte to read from (or thought of as n bytes consumed). The string read, and errors if any
func ReadString(data []byte) (int, string, error) {
	l := int(binary.LittleEndian.Uint32(data))
	data = data[4:]
	if len(data) < l {
		return 0, "", fmt.Errorf("not enough bytes to read string. Expected %d, but only have %d bytes", l, len(data))
	}
	return Align(l) + 4, string(data[:l]), nil
}

// Reads a Chrome-Pickle UTF-16 string.
// returns the next index of the passed []byte to read from (or thought of as n bytes consumed). The string read, and errors if any
func ReadString16(data []byte) (int, string, error) {
	l := int(binary.LittleEndian.Uint32(data)) // amount of []uint16, not []byte
	data = data[4:]
	if len(data) < l*2 {
		log.WithField("len(data)", len(data)).WithField("l", l).Error("not enough bytes for string16")
		return 0, "", errors.New("not enough bytes to read string16")
	}

	// convert []byte to []uint16 for utf16 lib conversion
	u16 := make([]uint16, l)
	for i := 0; i < l; i++ {
		u16[i] = binary.LittleEndian.Uint16(data[i*2:])
	}
	return Align(l*2) + 4, string(utf16.Decode(u16)), nil
}

// Reads Chrome-pickle length-prefixed bytes
// returns the next index of the passed []byte to read from (or thought of as n bytes consumed). The data read, and errors if any
func ReadBytes(data []byte) (int, []byte, error) {
	// @todo: copy bytes into new slice?
	// so the underlying array can be modified?
	l := int(binary.LittleEndian.Uint32(data))
	data = data[4:]
	if len(data) < l {
		return 0, nil, errors.New("not enough bytes to read")
	}
	return Align(l) + 4, data[:l], nil
}

// Align reads, pointers, to 32-bit boundary (4-byte)
func Align(n int) int { return n + (4-(n%4))%4 }
