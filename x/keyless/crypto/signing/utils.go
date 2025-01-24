package signing

// padTo32 pads a byte slice to 32 bytes. If the input is longer than 32 bytes,
// it returns the rightmost 32 bytes.
func padTo32(b []byte) []byte {
	if len(b) >= 32 {
		return b[len(b)-32:]
	}
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b)
	return padded
}
