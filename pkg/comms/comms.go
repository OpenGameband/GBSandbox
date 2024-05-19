package comms

import (
	"bytes"
	"errors"
	"github.com/sstallion/go-hid"
	"math"
	"time"
)

type Gameband struct {
	dev *hid.Device
}

func OpenHid() (*Gameband, error) {
	dev, err := hid.OpenFirst(0x2a90, 0x0021)
	return &Gameband{dev}, err
}

func (g *Gameband) Write(data []byte, bufSize uint, check uint8) ([]byte, bool, error) {
	n, err := g.dev.Write(data)
	if err != nil {
		return nil, false, err
	}
	if n <= 0 {
		return nil, false, nil
	}

	err = g.dev.SetNonblock(false)
	if err != nil {
		return nil, false, err
	}

	resp := make([]byte, 64)
	n, err = g.dev.ReadWithTimeout(resp, 5*time.Second)
	if n <= 0 {
		return nil, false, nil
	}

	outResp := resp[0:bufSize]

	if resp[0] != check {
		return nil, false, nil
	}
	if resp[1] != 0 {
		return nil, false, nil
	}
	return outResp, true, nil
}

func (g *Gameband) WriteAtOffset(commandCode uint, data []byte, offset, dataSize uint) error {
	if uint(len(data)) < offset+dataSize {
		return errors.New("error. Each data write must be of dataSize words")
	} else {
		buf := make([]byte, 37)
		buf[1] = 6
		buf[3] = byte(commandCode)
		buf[4] = byte(commandCode >> 8)

		for i := uint(0); i < dataSize; i++ {
			buf[i+5] = data[offset+i]
		}

		return g.BlindWrite(buf, 7)
	}
}

func (g *Gameband) ReadGameband() ([]byte, error) {
	buf := make([]byte, 0, 4096)

	offset := 6144
	for i := 0; i < 128; i++ {
		resp, err := g.ReadChunk(uint16(offset + i*16))
		if err != nil {
			return nil, err
		}
		buf = append(buf, resp...)
	}

	return buf, nil
}

func (g *Gameband) ReadChunk(offset uint16) ([]byte, error) {
	buf := []byte{0, 8, 0, 0, 0} // command (I think)
	buf[3] = byte(offset)
	buf[4] = byte(offset >> 8)

	if resp, good, err := g.Write(buf, 34, 9); good {
		return resp[2:], err
	}
	return nil, nil
}

func (g *Gameband) WriteData(data []byte) error {
	// Align the size of the array to a multiple of 32
	align := int(math.Ceil(float64(len(data))/128) * 128)
	if align != len(data) {
		paddingLen := align - len(data)
		padding := make([]byte, paddingLen)
		data = append(data, padding...)
	}

	if err := SanityCheck(data); err != nil {
		return err
	}

	err := g.SetDataLength(6144, uint16(128))
	if err != nil {
		return err
	}

	for i := 0; i < len(data); i += 32 {
		err := g.WriteAtOffset(6144+uint(i/2), data, uint(i), 32)
		if err != nil {
			return err
		}

		//resp, err := g.ReadChunk(uint16(6144 + (i)))
		//TODO: Verify data as it's written
	}

	return g.Commit() // Not sure what this does
}

func (g *Gameband) SetDataLength(offset, value uint16) error {
	data := make([]byte, 7)
	data[0] = 0
	data[1] = 4
	data[3] = byte(offset)
	data[4] = byte(offset >> 8)
	data[5] = byte(value)
	data[6] = byte(value >> 8)

	return g.BlindWrite(data, 5)
}

func (g *Gameband) BlindWrite(data []byte, check byte) error {
	_, ok, err := g.Write(data, 2, check)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("error sending command to Gameband")
	}
	return nil
}

func (g *Gameband) WritePayload(command uint16, data []byte, offset uint8, dataSize uint8) error {
	if len(data) < int(offset+dataSize) {
		return errors.New("payload must be of dataSize words")
	}

	buf := new(bytes.Buffer)
	buf.WriteByte(0)
	buf.WriteByte(6)

	buf.WriteByte(byte(command))
	buf.WriteByte(byte(command >> 8))

	for i := 0; i < int(dataSize); i++ {
		buf.WriteByte(data[int(offset)+i])
	}

	return g.BlindWrite(buf.Bytes(), 7)
}

func (g *Gameband) SetTime() error {
	encoded := make([]byte, 9)
	encoded[1] = 2
	PackTime(encoded, 5, uint32(time.Now().Unix()))
	return g.BlindWrite(encoded, 3)
}

func (g *Gameband) Commit() error {
	data := []byte{0, 10, 0}
	return g.BlindWrite(data, 11)
}
