package aptos

// StringToBCS prepend ULEB128 encoding of the string's byte length to the string's bytes
func StringToBCS(s string) []byte {
	v := []byte(s)
	return append(ULEB128Encode(len(v)), v...)
}
