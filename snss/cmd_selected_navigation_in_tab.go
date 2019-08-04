package snss

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type SelectNavigationInTab struct {
	TabID     int32
	Index     int32
	Timestamp int64
}

func (s SelectNavigationInTab) String() string {
	return fmt.Sprintf("Selected Navigation In Tab: Tab: %d, Index: %d, Time: %s\n", s.TabID, s.Index, timeConvert(s.Timestamp).Format("Jan 2 2006, 3:04:05 pm MST"))
}

func NewSelectNavigationInTab(data []byte) (SelectNavigationInTab, error) {
	if len(data) < 16 {
		return SelectNavigationInTab{}, errors.New("not enough data to create SelectNavInTab")
	}

	return SelectNavigationInTab{
		TabID:     int32(binary.LittleEndian.Uint32(data)),
		Index:     int32(binary.LittleEndian.Uint32(data[4:])),
		Timestamp: int64(binary.LittleEndian.Uint64(data[8:])),
	}, nil
}
