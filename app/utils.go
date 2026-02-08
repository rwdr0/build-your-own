package main

import "os"

func GetArgumentsForStage(indexes ...int) []string {
	var args []string
	for _, idx := range indexes {
		if idx < len(os.Args) {
			args = append(args, os.Args[idx])
		}
	}
	return args
}
