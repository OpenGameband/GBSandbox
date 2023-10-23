package comms

func checksum(data []byte) (uint16, uint16) {
	var cs1, cs2 uint32
	buf := data
	dataLen := len(data)

	for i := 0; i < dataLen; i++ {
		dataPiece := buf[i]

		cs1 = (cs1 + (uint32(dataPiece) & 255)) % 255
		cs2 = (cs2 + cs1) % 255
	}
	return uint16(cs1), uint16(cs2)
}

func PackTime(buf []byte, offset uint, seconds uint32) {
	for i := uint(0); i < 4; i++ {
		buf[offset+i] = byte(seconds >> (i * 8))
	}
}

func PackUInt32(seconds uint32) []byte {
	buf := make([]byte, 8)
	for i := uint(0); i < 8; i += 2 {
		buf[i] = byte(seconds >> ((i / 2) * 8))
	}
	return buf
}

func PackUInt16(val uint16) []byte {
	buf := make([]byte, 2)
	for i := uint(0); i < 2; i++ {
		buf[i] = byte(val >> (i * 8))
	}
	return buf
}
