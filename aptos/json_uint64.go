package aptos

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fardream/go-bcs/bcs"
)

// JsonUint64 is an uint64, but serialized into a string, and can be deserialized from either a string or a number from json.
// This is because aptos fullnode uses string for uint64, whereas golang's json encoding only support number.
type JsonUint64 uint64

var (
	_ json.Marshaler   = (*JsonUint64)(nil)
	_ json.Unmarshaler = (*JsonUint64)(nil)
	_ bcs.Marshaler    = (*JsonUint64)(nil)
)

func (i *JsonUint64) UnmarshalJSON(data []byte) error {
	var j uint64
	if err := json.Unmarshal(data, &j); err != nil {
		var str string
		if errStr := json.Unmarshal(data, &str); errStr != nil {
			return fmt.Errorf("failed to parse as uint64: %w, then failed to parse as string: %v", err, errStr)
		}
		if j, err := strconv.ParseUint(str, 10, 64); err != nil {
			return fmt.Errorf("failed to parse string: %s as uint64: %w", str, err)
		} else {
			*i = JsonUint64(j)
		}
	} else {
		*i = JsonUint64(j)
	}

	return nil
}

func (i JsonUint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(i), 10))
}

// MarshalBCS marshals the uint64 to bcs bytes.
func (i JsonUint64) MarshalBCS() ([]byte, error) {
	return bcs.Marshal(uint64(i))
}

// ToBCS calls bcs.Marshal and returns the bytes, ignoring the error if there is any.
//
// Deprecated: Use [JsonUint64.MarshalBCS] directly
func (i JsonUint64) ToBCS() []byte {
	r, _ := bcs.Marshal(uint64(i))

	return r
}
