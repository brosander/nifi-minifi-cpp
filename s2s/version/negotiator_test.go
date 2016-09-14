package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func create(t *testing.T) *StandardVersionNegotiator {
	s, err := NewStandardVersionNegotiator([]uint8{5, 4, 3, 2})
	assert.Nil(t, err)
	assert.Equal(t, []uint8{5, 4, 3, 2}, s.SupportedVersions())
	return s
}

func TestConstructor(t *testing.T) {
	create(t)
}

func TestPreferredVersion(t *testing.T) {
	s := create(t)
	assert.Equal(t, uint8(5), s.PreferredVersion())
}

func TestVersion(t *testing.T) {
	s := create(t)
	assert.Equal(t, uint8(5), s.Version())
}

func TestSetVersion(t *testing.T) {
	s := create(t)
	assert.Nil(t, s.SetVersion(4))
	assert.NotNil(t, s.SetVersion(6))
	assert.Equal(t, uint8(4), s.Version())
}

func TestIsVersionSupported(t *testing.T) {
	s := create(t)
	assert.False(t, s.IsVersionSupported(6))
	assert.True(t, s.IsVersionSupported(5))
}

func TestPreferredVersionWithMax(t *testing.T) {
	s := create(t)
	_, err := s.PreferredVersionWithMax(1)
	assert.NotNil(t, err)

	v, err := s.PreferredVersionWithMax(3)
	assert.Equal(t, uint8(3), v)
	assert.Nil(t, err)
}
