package utils

func RlpNextItemSize(data []byte) int {
	if len(data) == 0 {
		return -1
	}

	prefix := data[0]

	switch {
	case prefix <= 0x7f:
		// Single byte
		return 1

	case prefix <= 0xb7:
		// Short string
		return int(prefix - 0x80 + 1)

	case prefix <= 0xbf:
		// Long string
		lengthSize := int(prefix - 0xb7)
		if len(data) < lengthSize+1 {
			return -1
		}
		length := int(data[1])
		for i := 2; i < lengthSize+1; i++ {
			length = (length << 8) + int(data[i])
		}
		return length + lengthSize + 1

	case prefix <= 0xf7:
		// Short list
		return int(prefix - 0xc0 + 1)

	case prefix <= 0xff:
		// Long list
		lengthSize := int(prefix - 0xf7)
		if len(data) < lengthSize+1 {
			return -1
		}
		length := int(data[1])
		for i := 2; i < lengthSize+1; i++ {
			length = (length << 8) + int(data[i])
		}
		return length + lengthSize + 1

	default:
		return -1
	}
}
