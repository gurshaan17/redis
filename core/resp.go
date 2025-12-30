package core

import (
	"errors"
)

// Decode decodes a RESP-encoded value from bytes
func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}

	value, _, err := DecodeOne(data)
	return value, err
}

// DecodeOne decodes a single RESP value and returns:
// value, bytes consumed (delta), error
func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	default:
		return nil, 0, errors.New("unknown RESP type")
	}
}

// +OK\r\n
func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}
	if pos+1 >= len(data) {
		return "", 0, errors.New("invalid simple string")
	}
	return string(data[1:pos]), pos + 2, nil
}

// -ERR something\r\n
func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

// :123\r\n
func readInt64(data []byte) (int64, int, error) {
	pos := 1
	var value int64

	for pos < len(data) && data[pos] != '\r' {
		value = value*10 + int64(data[pos]-'0')
		pos++
	}
	if pos+1 >= len(data) {
		return 0, 0, errors.New("invalid integer")
	}
	return value, pos + 2, nil
}

// $5\r\nhello\r\n
func readBulkString(data []byte) (string, int, error) {
	pos := 1

	length, delta := readLength(data[pos:])
	pos += delta

	if length == -1 {
		return "", pos, nil
	}

	end := pos + length
	if end+2 > len(data) {
		return "", 0, errors.New("invalid bulk string")
	}

	return string(data[pos:end]), end + 2, nil
}

// reads length until CRLF
func readLength(data []byte) (int, int) {
	pos := 0
	length := 0
	sign := 1

	if data[pos] == '-' {
		sign = -1
		pos++
	}

	for pos < len(data) {
		b := data[pos]
		if b == '\r' {
			return sign * length, pos + 2
		}
		length = length*10 + int(b-'0')
		pos++
	}
	return 0, 0
}

// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
func readArray(data []byte) ([]interface{}, int, error) {
	pos := 1

	count, delta := readLength(data[pos:])
	pos += delta

	if count == -1 {
		return nil, pos, nil
	}

	elems := make([]interface{}, count)

	for i := 0; i < count; i++ {
		elem, d, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		pos += d
	}

	return elems, pos, nil
}
