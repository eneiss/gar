package utils

// import "math"
// nah we're doing bit shifts instead

// Takes an array of bytes representing a base 8 number, and returns
// the base 10 equivalent as an integer.
// Does NOT check for overflow with bit shifts!!
func Base8ToBase10(bytes []byte) int {
	res := 0

	length := len(bytes)

	for i, b := range bytes {
		// bit shift shenanigans to multiply by powers of 8:
		// multiplying by 8 means shifting bits to the left 3 times, without overflow considerations
		// and substracting 0x30 to go from ASCII code to raw integer value
		res += int(b-byte(0x30)) << (3 * (length - i - 1))

		// easier to understand, worse efficiency
		// res += (int(b) - 0x30) * int(math.Pow(float64(8), float64(length-i-1)))
	}

	return res
}
