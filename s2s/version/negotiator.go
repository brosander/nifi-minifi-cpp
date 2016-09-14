package version

import (
	"errors"
	"fmt"
)

type VersionNegotiator interface {
	Version() uint8
	SetVersion(version uint8) error
	PreferredVersion() uint8
	PreferredVersionWithMax(max uint8) (uint8, error)
	IsVersionSupported(version uint8) bool
	SupportedVersions() []uint8
}

type StandardVersionNegotiator struct {
	version           uint8
	supportedVersions []uint8
}

func copyUints(val []uint8) []uint8 {
	result := make([]uint8, len(val))
	copy(result, val)
	return result
}

func NewStandardVersionNegotiator(supportedVersions []uint8) (*StandardVersionNegotiator, error) {
	if supportedVersions == nil || len(supportedVersions) == 0 {
		return nil, errors.New("At least one version must be supported")
	}
	return &StandardVersionNegotiator{supportedVersions[0], copyUints(supportedVersions)}, nil
}

func (s *StandardVersionNegotiator) Version() uint8 {
	return s.version
}

func (s *StandardVersionNegotiator) SetVersion(version uint8) error {
	if s.IsVersionSupported(version) {
		s.version = version
		return nil
	}
	return fmt.Errorf("Version %d not in supported versions (%v)", version, s.supportedVersions)
}

func (s *StandardVersionNegotiator) PreferredVersion() uint8 {
	return s.supportedVersions[0]
}

func (s *StandardVersionNegotiator) IsVersionSupported(version uint8) bool {
	for _, v := range s.supportedVersions {
		if v == version {
			return true
		}
	}
	return false
}

func (s *StandardVersionNegotiator) PreferredVersionWithMax(max uint8) (uint8, error) {
	for _, v := range s.supportedVersions {
		if v <= max {
			return v, nil
		}
	}
	return 0, fmt.Errorf("Unable to find supporte version <= %d in %v", max, s.supportedVersions)
}

func (s *StandardVersionNegotiator) SupportedVersions() []uint8 {
	return copyUints(s.supportedVersions)
}

// Ensure StandardVersionNegotiator implements VersionNegotiator
var _ VersionNegotiator = (*StandardVersionNegotiator)(nil)
