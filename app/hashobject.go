package main

import "github.com/rwdr0/build-your-own/git/app/utils"

func HashObject() {
	sourcePath := utils.GetArgumentsForStage(3)[0]
	utils.HashObject(sourcePath, true)
}
