// Code generated by "stringer -type AuxClobMarketOrderStatus -linecomment"; DO NOT EDIT.

package aptos

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AuxClobMarketOrderStatus_Placed-0]
	_ = x[AuxClobMarketOrderStatus_Filled-1]
	_ = x[AuxClobMarketOrderStatus_Cancelled-2]
}

const _AuxClobMarketOrderStatus_name = "PlacedFilledCancelled"

var _AuxClobMarketOrderStatus_index = [...]uint8{0, 6, 12, 21}

func (i AuxClobMarketOrderStatus) String() string {
	if i >= AuxClobMarketOrderStatus(len(_AuxClobMarketOrderStatus_index)-1) {
		return "AuxClobMarketOrderStatus(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _AuxClobMarketOrderStatus_name[_AuxClobMarketOrderStatus_index[i]:_AuxClobMarketOrderStatus_index[i+1]]
}
