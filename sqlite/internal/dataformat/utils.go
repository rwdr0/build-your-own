package dataformat

func GetContentSize(serialType uint64) uint64 {
	switch {
	case serialType == 1:
		return 1 // 8-bit integer
	case serialType == 2:
		return 2 // 16-bit integer
	case serialType == 3:
		return 3 // 24-bit integer
	case serialType == 4:
		return 4 // 32-bit integer
	case serialType == 5:
		return 6 // 48-bit integer
	case serialType == 6:
		return 8 // 64-bit integer
	case serialType == 7:
		return 8 // 64-bit floating point
	case serialType >= 12 && serialType%2 == 0:
		return (serialType - 12) / 2 // BLOB
	case serialType >= 13 && serialType%2 == 1:
		return (serialType - 13) / 2 // Text string
	default:
		return 0 // NULL
	}
}
