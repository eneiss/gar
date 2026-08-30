package utils

import (
	"fmt"
	// "math"
	// nah we're doing bit shifts instead
)

// Takes an array of bytes representing a base 8 number, and returns
// the base 10 equivalent as an integer.
// Does NOT check for overflow with bit shifts, uses uint64 to guard against
// most cases instead.
func Base8ToBase10(bytes []byte) (uint64, error) {
	var res uint64 = 0

	length := len(bytes)

	for i, b := range bytes {

		if b < '0' || b > '7' {
			return 0, fmt.Errorf("invalid base 8 character: '%c' (0x%x)", b, b)
		}
		// bit shift shenanigans to multiply by powers of 8:
		// multiplying by 8 means shifting bits to the left 3 times, without overflow considerations
		// and substracting 0x30 to go from ASCII code to raw integer value
		res += uint64(b-'0') << (3 * (length - i - 1))
	}

	return res, nil
}
