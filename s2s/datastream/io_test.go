package datastream

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DIO struct {
	DataInput
	DataOutput
	buf *bytes.Buffer
}

func NewDIO() *DIO {
	buf := new(bytes.Buffer)
	return &DIO{DataInput{buf}, DataOutput{buf}, buf}
}

func TestBool(t *testing.T) {
	roundTrip := []bool{true, false}
	for _, tt := range roundTrip {
		dio := NewDIO()
		assert.Nil(t, dio.WriteBool(tt))
		b, err := dio.ReadBool()
		assert.Nil(t, err)
		assert.Equal(t, tt, b)
	}
}

func TestSignedByte(t *testing.T) {
	roundTrip := []int8{math.MinInt8, -1, 0, 1, math.MaxInt8}
	for _, tt := range roundTrip {
		dio := NewDIO()
		assert.Nil(t, dio.WriteSignedByte(tt))
		b, err := dio.ReadSignedByte()
		assert.Nil(t, err)
		assert.Equal(t, tt, b)
	}
}

func TestByte(t *testing.T) {
	roundTrip := []byte{0, 1, 127, math.MaxUint8}
	for _, tt := range roundTrip {
		dio := NewDIO()
		assert.Nil(t, dio.WriteByte(tt))
		b, err := dio.ReadByte()
		assert.Nil(t, err)
		assert.Equal(t, tt, b)
	}
}

func TestChar(t *testing.T) {
	roundTrip := []rune{0, 1, 127, 255, '\u2318'}
	for _, tt := range roundTrip {
		dio := NewDIO()
		assert.Nil(t, dio.WriteChar(tt))
		b, err := dio.ReadChar()
		assert.Nil(t, err)
		assert.Equal(t, tt, b)
	}
}

func TestDouble(t *testing.T) {
	roundTrip := []float64{0, 1, 127, 255, 0.125, math.MaxFloat64, math.SmallestNonzeroFloat64}
	for _, tt := range roundTrip {
		dio := NewDIO()
		assert.Nil(t, dio.WriteDouble(tt))
		b, err := dio.ReadDouble()
		assert.Nil(t, err)
		assert.Equal(t, tt, b)
	}
}

func TestFloat(t *testing.T) {
	roundTrip := []float32{0, 1, 127, 255, 0.125, math.MaxFloat32, math.SmallestNonzeroFloat32}
	for _, tt := range roundTrip {
		dio := NewDIO()
		assert.Nil(t, dio.WriteFloat(tt))
		b, err := dio.ReadFloat()
		assert.Nil(t, err)
		assert.Equal(t, tt, b)
	}
}

func TestInt(t *testing.T) {
	roundTrip := []int32{math.MinInt32, -1, 0, 1, 127, 255, math.MaxInt32}
	for _, tt := range roundTrip {
		dio := NewDIO()
		assert.Nil(t, dio.WriteInt(tt))
		b, err := dio.ReadInt()
		assert.Nil(t, err)
		assert.Equal(t, tt, b)
	}
}

func TestShort(t *testing.T) {
	roundTrip := []int16{math.MinInt16, -1, 0, 1, 127, 255, math.MaxInt16}
	for _, tt := range roundTrip {
		dio := NewDIO()
		assert.Nil(t, dio.WriteShort(tt))
		b, err := dio.ReadShort()
		assert.Nil(t, err)
		assert.Equal(t, tt, b)
	}
}

func TestUnsignedShort(t *testing.T) {
	roundTrip := []uint16{0, 1, 127, 255, math.MaxUint16}
	for _, tt := range roundTrip {
		dio := NewDIO()
		assert.Nil(t, dio.WriteUnsignedShort(tt))
		b, err := dio.ReadUnsignedShort()
		assert.Nil(t, err)
		assert.Equal(t, tt, b)
	}
}

func TestUTF(t *testing.T) {
	roundTrip := []string{"The quick brown fox", "j\u2318mped over the lazy dog", "\u0000null chars\u0000", "\u0722two byte unicode\u0723"}
	for _, tt := range roundTrip {
		dio := NewDIO()
		assert.Nil(t, dio.WriteUTF(tt))
		b, err := dio.ReadUTF()
		assert.Nil(t, err)
		assert.Equal(t, tt, b)
	}
}
