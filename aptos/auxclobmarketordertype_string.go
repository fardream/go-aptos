// Code generated by "stringer -type AuxClobMarketOrderType -linecomment"; DO NOT EDIT.

package aptos

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AuxClobMarketOrderType_Limit-100]
	_ = x[AuxClobMarketOrderType_FOK-101]
	_ = x[AuxClobMarketOrderType_IOC-102]
	_ = x[AuxClobMarketOrderType_POST_ONLY-103]
	_ = x[AuxClobMarketOrderType_PASSIVE_JOIN-104]
}

const _AuxClobMarketOrderType_name = "LIMITFOKIOCPOST_ONLYPASSIVE_JOIN"

var _AuxClobMarketOrderType_index = [...]uint8{0, 5, 8, 11, 20, 32}

func (i AuxClobMarketOrderType) String() string {
	i -= 100
	if i >= AuxClobMarketOrderType(len(_AuxClobMarketOrderType_index)-1) {
		return "AuxClobMarketOrderType(" + strconv.FormatInt(int64(i+100), 10) + ")"
	}
	return _AuxClobMarketOrderType_name[_AuxClobMarketOrderType_index[i]:_AuxClobMarketOrderType_index[i+1]]
}