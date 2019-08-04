package snss

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

const MAGICBYTES = 0x53534E53                // 'SNSS'
const EPOCHCONVERT int64 = 11644473600000000 // convert 'FILETIME' (see link) to unix date
// https://docs.microsoft.com/en-us/windows/win32/api/minwinbase/ns-minwinbase-filetime

type SNSSFileHeader struct {
	Signature int32
	Version   int32
}

type CommandID int

const (
	CommandUpdateTabNavigation CommandID = iota + 1
	CommandRestoredEntry
	CommandWindow
	CommandSelectedNavigationInTab
	CommandPinnedState
	CommandSetExtensionAppID
)

func (c CommandID) String() string {
	names := [...]string{
		"Invalid",
		"Update Tab Navigation",
		"Restored Entry",
		"Window",
		"Selected Navigation In Tab",
		"Pinned State",
		"Set Extension App ID",
	}
	if c > CommandSetExtensionAppID || c < CommandUpdateTabNavigation {
		return "Unknown"
	}
	return names[c]
}

/* returns file version on success, otherwise error */
func FileInfo(f io.Reader) (int, error) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(f, buf)
	if err != nil {
		log.WithError(err).Error("error reading file header")
		return -1, err
	}

	head := SNSSFileHeader{
		Signature: int32(binary.LittleEndian.Uint32(buf)),
		Version:   int32(binary.LittleEndian.Uint32(buf[4:])),
	}

	if head.Signature != MAGICBYTES {
		log.WithField("got", head.Signature).WithField("expected", MAGICBYTES).Error("unexpected magic bytes")
		return 0, errors.New(fmt.Sprintf("error parsing SNSS file. Invalid magic bytes: 0x%x\n", head.Signature))
	}
	return int(head.Version), nil
}

func ReadCommand(f io.Reader) error {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(f, buf); err != nil {
		if err != io.EOF {
			log.WithError(err).Error("command size read failure")
		}
		return err
	}
	size := binary.LittleEndian.Uint16(buf)
	log.WithField("size", size).Debug("command data size determined")

	data := make([]byte, size)
	if _, err := io.ReadFull(f, data); err != nil {
		log.WithError(err).Error("command data read failure")
		return err
	}

	cid, data := int(data[0]), data[1:]
	cmd := CommandID(cid)
	log.WithField("command", cmd).WithField("id", cid).Debug("command determined")

	switch cmd {
	case CommandSelectedNavigationInTab:
		c, err := NewSelectNavigationInTab(data)
		if err != nil {
			log.WithError(err).Error("error creating selected nav tab")
			return err
		}
		fmt.Printf("%s\n", c)
	case CommandUpdateTabNavigation:
		c, err := NewUpdateTabNavigation(data)
		if err != nil {
			log.WithError(err).Error("error creation update tab nav cmd")
			return err
		}
		fmt.Printf("%+v\n", c)
	}

	//fmt.Printf("--%v--\n", cmd)
	fmt.Println("----")
	return nil
}

func timeConvert(t int64) time.Time {
	return time.Unix((t-EPOCHCONVERT)/1000000, 0)
}
