package aptos

func StringToBCS(s string) []byte {
	v := []byte(s)
	return append(ULEB128Encode(len(v)), v...)
}
