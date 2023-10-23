package comms

import (
	"bytes"
	"errors"
)

type GamebandHeader struct {
	Timezone                uint8
	AltTimezone             uint8
	TzChange                uint32
	Orientation             uint8
	TransitionFrameDuration uint8
	ScreenCount             uint8
	AnimationDataLength     uint16
	Checksum0               uint8
	Checksum1               uint8
}

// AnimationHeader contains data about the screen
type AnimationHeader struct {
	ScreenType    uint8
	PauseMode     uint8
	PauseDuration uint16
	FrameDuration uint16
	AnimationType uint8
	DataLength    uint16
}

type Animation struct {
	Header AnimationHeader
	Frames []Frame
}

type Frame struct {
	Data [20][7]bool // the rows/cols of the gameband screen
}

type GBData struct {
	Header     GamebandHeader
	Animations []Animation
}

func PackGamebandHeader(header GamebandHeader) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(header.Timezone)
	buf.WriteByte(0)
	buf.WriteByte(header.AltTimezone)
	buf.WriteByte(0)
	buf.Write(PackUInt32(header.TzChange))
	buf.WriteByte(header.Orientation)
	buf.WriteByte(0)
	buf.WriteByte(header.TransitionFrameDuration)
	buf.WriteByte(0)
	buf.WriteByte(header.ScreenCount)
	buf.WriteByte(0)
	buf.Write(PackUInt16(header.AnimationDataLength))

	return buf.Bytes()
}

func SanityCheck(data []byte) error {
	if len(data) > 2048 {
		return errors.New("data is too large")
	}
	// the original software also checked the signed nature of the variables, we have uints so it's fine
	return nil
}
