package utils

import "os"

// GetArgumentsForStage extracts os.Args values at the given positional indexes.
// Each index corresponds to a position in the command-line argument list.
func GetArgumentsForStage(indexes ...int) []string {
	var args []string
	for _, idx := range indexes {
		if idx < len(os.Args) {
			args = append(args, os.Args[idx])
		}
	}
	return args
}
