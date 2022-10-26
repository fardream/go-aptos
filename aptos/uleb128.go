package aptos

import "fmt"

// ULEB128Encode converts an integer into []byte
func ULEB128Encode[T uint8 | uint16 | uint32 | uint64 | uint | int](input T) []byte {
	var result []byte

	for {
		b := (byte)(input & 127)
		input >>= 7

		if input == 0 {
			result = append(result, b)
			break
		} else {
			result = append(result, b|128)
		}
	}

	return result
}

// ULEB128Decode decodes byte array into an integer, returns the decoded value, the number of bytes consumed, and a possible error.
// If error is returned, the number of bytes returned is guaranteed to be 0.
func ULEB128Decode[T uint8 | uint16 | uint32 | uint64 | uint | int](data []byte) (T, int, error) {
	var v, shift T
	for i := 0; i < len(data); i++ {
		d := T(data[i])
		ld := d & 127
		if (ld<<shift)>>shift != ld {
			return v, 0, fmt.Errorf("overflow at index %d: %v", i, ld)
		}
		ld <<= shift
		v = ld + v
		if v < ld {
			return v, 0, fmt.Errorf("overflow after adding index %d: %v %v", i, ld, v)
		}
		if d < 128 {
			return v, i + 1, nil
		}
		shift += 7
	}

	return v, 0, fmt.Errorf("failed to find the highest significant 7 bits: %v", v)
}
