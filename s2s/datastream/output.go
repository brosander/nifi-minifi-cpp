package datastream

import (
	"encoding/binary"
	"errors"
	"io"
	"unicode/utf16"
)

type DataOutput struct {
	writer io.Writer
}

func (d *DataOutput) WriteBool(val bool) error {
	if val {
		return d.WriteByte(1)
	}
	return d.WriteByte(0)
}

func (d *DataOutput) WriteByte(val byte) error {
	_, err := d.writer.Write([]byte{val})
	return err
}

func (d *DataOutput) WriteSignedByte(val int8) error {
	return binary.Write(d.writer, binary.BigEndian, val)
}

func (d *DataOutput) WriteChar(val rune) error {
	return binary.Write(d.writer, binary.BigEndian, utf16.Encode([]rune{val})[0])
}

func (d *DataOutput) WriteDouble(val float64) error {
	return binary.Write(d.writer, binary.BigEndian, val)
}

func (d *DataOutput) WriteFloat(val float32) error {
	return binary.Write(d.writer, binary.BigEndian, val)
}

func (d *DataOutput) WriteFully(val []byte) error {
	_, err := d.writer.Write(val)
	return err
}

func (d *DataOutput) WriteInt(val int32) error {
	return binary.Write(d.writer, binary.BigEndian, val)
}

func (d *DataOutput) WriteLong(val int64) error {
	return binary.Write(d.writer, binary.BigEndian, val)
}

func (d *DataOutput) WriteShort(val int16) error {
	return binary.Write(d.writer, binary.BigEndian, val)
}

func (d *DataOutput) WriteUnsignedShort(val uint16) error {
	return binary.Write(d.writer, binary.BigEndian, val)
}

func (d *DataOutput) WriteUTF(val string) error {
	if len(val) == 0 {
		return nil
	}
	bytes := make([]byte, 0, len(val))
	for _, r := range val {
		if r > 0x0001 && r < 0x007F {
			bytes = append(bytes, byte(r))
		} else if r > 0x07FF && r <= 0xFFFF {
			bytes = append(bytes, byte(0xE0|((r>>12)&0x0F)), byte(0x80|((r>>6)&0x3F)), byte(0x80|r&0x3F))
		} else if r < 0x0800 {
			bytes = append(bytes, byte(0xC0|((r>>6)&0x1F)), byte(0x80|r&0x3F))
		} else {
			return errors.New("Can't write values greater than 0xFFFF")
		}
	}
	d.WriteUnsignedShort(uint16(len(bytes)))
	_, err := d.writer.Write(bytes)
	return err
}
