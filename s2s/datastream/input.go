package datastream

import (
	"encoding/binary"
	"errors"
	"io"
	"unicode/utf16"
)

type DataInput struct {
	reader io.Reader
}

func (d *DataInput) ReadBool() (bool, error) {
	b, err := d.ReadByte()
	if err != nil {
		return false, err
	}
	return b != 0, nil
}

func (d *DataInput) ReadSignedByte() (int8, error) {
	var result int8
	if err := binary.Read(d.reader, binary.BigEndian, &result); err != nil {
		return 0, err
	}
	return result, nil
}

func (d *DataInput) ReadByte() (byte, error) {
	var result byte
	if err := binary.Read(d.reader, binary.BigEndian, &result); err != nil {
		return 0, err
	}
	return result, nil
}

func (d *DataInput) ReadChar() (rune, error) {
	var result uint16
	if err := binary.Read(d.reader, binary.BigEndian, &result); err != nil {
		return 0, err
	}
	return utf16.Decode([]uint16{result})[0], nil
}

func (d *DataInput) ReadDouble() (float64, error) {
	var result float64
	if err := binary.Read(d.reader, binary.BigEndian, &result); err != nil {
		return 0, err
	}
	return result, nil
}

func (d *DataInput) ReadFloat() (float32, error) {
	var result float32
	if err := binary.Read(d.reader, binary.BigEndian, &result); err != nil {
		return 0, err
	}
	return result, nil
}

func (d *DataInput) ReadFully(buf []byte) error {
	if err := binary.Read(d.reader, binary.BigEndian, &buf); err != nil {
		return err
	}
	return nil
}

func (d *DataInput) ReadInt() (int32, error) {
	var result int32
	if err := binary.Read(d.reader, binary.BigEndian, &result); err != nil {
		return 0, err
	}
	return result, nil
}

func (d *DataInput) ReadLong() (int64, error) {
	var result int64
	if err := binary.Read(d.reader, binary.BigEndian, &result); err != nil {
		return 0, err
	}
	return result, nil
}

func (d *DataInput) ReadShort() (int16, error) {
	var result int16
	if err := binary.Read(d.reader, binary.BigEndian, &result); err != nil {
		return 0, err
	}
	return result, nil
}

func (d *DataInput) ReadUnsignedShort() (uint16, error) {
	var result uint16
	if err := binary.Read(d.reader, binary.BigEndian, &result); err != nil {
		return 0, err
	}
	return result, nil
}

func (d *DataInput) ReadUTF() (string, error) {
	var numBytes uint16
	if err := binary.Read(d.reader, binary.BigEndian, &numBytes); err != nil {
		return "", err
	}
	bytes := make([]byte, numBytes)
	if err := binary.Read(d.reader, binary.BigEndian, &bytes); err != nil {
		return "", err
	}
	signedNumBytes := int(numBytes)
	decodedBytes := make([]byte, 0, numBytes)
	for count := 0; count < signedNumBytes; count++ {
		c := bytes[count]
		cShift := c >> 4
		switch {
		case cShift >= 0 && cShift < 8:
			decodedBytes = append(decodedBytes, c)
		case cShift > 11 && cShift <= 14:
			if signedNumBytes <= count+1 {
				return "", errors.New("Expected 2 bytes for code points starting with 110 or 10")
			}
			c2 := bytes[count+1]
			if (c2 & 0xC0) != 0x80 {
				return "", errors.New("Expected byte2 & 0xC0 to equal 0x80")
			}
			if cShift < 14 {
				if c == 0xC0 && c2 == 0x80 {
					decodedBytes = append(decodedBytes, 0x00)
				} else {
					decodedBytes = append(decodedBytes, c, c2)
				}
			} else {
				if signedNumBytes <= count+2 {
					return "", errors.New("Expected 3 bytes for code points starting with 1110")
				}
				c3 := bytes[count+2]
				if (c3 & 0xC0) != 0x80 {
					return "", errors.New("Expected byte3 & 0xC0 to equal 0x80")
				}
				decodedBytes = append(decodedBytes, c, c2, c3)
				count++
			}
			count++
		default:
			return "", errors.New("Did not expect to encounter 1111 xxxx byte")
		}
	}
	return string(decodedBytes), nil
}
