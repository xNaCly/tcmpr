package v2

// bits returns the bits in b as uint8
func bits(b byte) [8]uint8 {
	c := [8]uint8{}
	for i := range 8 {
		c[7-i] = (b >> i) & 1
	}
	return c
}
