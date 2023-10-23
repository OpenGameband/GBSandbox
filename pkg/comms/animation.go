package comms

import (
	"bytes"
	"fmt"
)

func packAnimationHeader(header AnimationHeader) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(header.ScreenType)
	buf.WriteByte(0)
	buf.WriteByte(header.PauseMode)
	buf.WriteByte(0)

	buf.WriteByte(byte(header.PauseDuration))
	buf.WriteByte(byte(header.PauseDuration >> 8))

	buf.WriteByte(byte(header.FrameDuration))
	buf.WriteByte(byte(header.FrameDuration >> 8))

	buf.WriteByte(header.AnimationType)
	buf.WriteByte(0)

	buf.WriteByte(byte(header.DataLength))
	buf.WriteByte(byte(header.DataLength >> 8))

	return buf.Bytes()
}

func (g *Gameband) WriteGBData(data GBData) error {
	buf := new(bytes.Buffer)
	buf.Write(PackGamebandHeader(data.Header))

	animationBuf := new(bytes.Buffer)

	for _, animation := range data.Animations {
		animationBuf.Write(packAnimationHeader(animation.Header))
		//TODO: Pack animation data too
	}

	cs0, cs1 := checksum(animationBuf.Bytes())
	fmt.Println("Checksum 0: ", cs0)
	fmt.Println("Checksum 1: ", cs1)

	buf.WriteByte(byte(cs0))
	buf.WriteByte(0)
	buf.WriteByte(byte(cs1))
	buf.WriteByte(0)

	buf.Write(animationBuf.Bytes())

	err := g.WriteData(buf.Bytes())
	return err
}
