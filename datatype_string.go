// Code generated by "stringer -type=DataType"; DO NOT EDIT.

package swlib

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DataTypeUnspec-0]
	_ = x[DataTypeInt-1]
	_ = x[DataTypeString-2]
	_ = x[DataTypePorts-3]
	_ = x[DataTypeLink-4]
	_ = x[DataTypeNoVal-5]
}

const _DataType_name = "DataTypeUnspecDataTypeIntDataTypeStringDataTypePortsDataTypeLinkDataTypeNoVal"

var _DataType_index = [...]uint8{0, 14, 25, 39, 52, 64, 77}

func (i DataType) String() string {
	if i >= DataType(len(_DataType_index)-1) {
		return "DataType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DataType_name[_DataType_index[i]:_DataType_index[i+1]]
}
