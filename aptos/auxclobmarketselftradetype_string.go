// Code generated by "stringer -type AuxClobMarketSelfTradeType -linecomment"; DO NOT EDIT.

package aptos

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AuxClobMarketSelfTradeType_CancelPassive-200]
	_ = x[AuxClobMarketSelfTradeType_CancelAggressive-201]
	_ = x[AuxClobMarketSelfTradeType_CancelBoth-202]
}

const _AuxClobMarketSelfTradeType_name = "CANCEL_PASSIVECANCEL_AGGRESSIVECANCEL_BOTH"

var _AuxClobMarketSelfTradeType_index = [...]uint8{0, 14, 31, 42}

func (i AuxClobMarketSelfTradeType) String() string {
	i -= 200
	if i >= AuxClobMarketSelfTradeType(len(_AuxClobMarketSelfTradeType_index)-1) {
		return "AuxClobMarketSelfTradeType(" + strconv.FormatInt(int64(i+200), 10) + ")"
	}
	return _AuxClobMarketSelfTradeType_name[_AuxClobMarketSelfTradeType_index[i]:_AuxClobMarketSelfTradeType_index[i+1]]
}
