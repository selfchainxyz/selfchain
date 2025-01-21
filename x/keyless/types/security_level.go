package types

// SecurityLevel represents the security level for key generation
type SecurityLevel int32

const (
	// SecurityLevel_UNSPECIFIED represents an unspecified security level
	SecurityLevel_UNSPECIFIED SecurityLevel = 0
	// SecurityLevel_LOW represents a low security level
	SecurityLevel_LOW SecurityLevel = 1
	// SecurityLevel_MEDIUM represents a medium security level
	SecurityLevel_MEDIUM SecurityLevel = 2
	// SecurityLevel_HIGH represents a high security level
	SecurityLevel_HIGH SecurityLevel = 3
)

// String returns the string representation of SecurityLevel
func (s SecurityLevel) String() string {
	switch s {
	case SecurityLevel_LOW:
		return "LOW"
	case SecurityLevel_MEDIUM:
		return "MEDIUM"
	case SecurityLevel_HIGH:
		return "HIGH"
	default:
		return "UNSPECIFIED"
	}
}

// IsValid returns whether the security level is valid
func (s SecurityLevel) IsValid() bool {
	switch s {
	case SecurityLevel_LOW, SecurityLevel_MEDIUM, SecurityLevel_HIGH:
		return true
	default:
		return false
	}
}
