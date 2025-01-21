package types

// IsValid returns whether the security level is valid
func (s SecurityLevel) IsValid() bool {
	switch s {
	case SecurityLevel_SECURITY_LEVEL_STANDARD,
		SecurityLevel_SECURITY_LEVEL_HIGH,
		SecurityLevel_SECURITY_LEVEL_ENTERPRISE:
		return true
	default:
		return false
	}
}
