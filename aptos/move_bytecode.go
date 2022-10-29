package aptos

import "encoding/json"

type MoveBytecode []byte

func (b MoveBytecode) String() string {
	return prefixedHexString(b)
}

func (b MoveBytecode) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *MoveBytecode) UnmarshalJSON(date []byte) error {
	var str string
	err := json.Unmarshal(date, &str)
	if err != nil {
		return nil
	}

	allBytes, err := parseHexString(str)
	if err != nil {
		return err
	}

	*b = allBytes

	return nil
}
