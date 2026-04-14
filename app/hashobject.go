package main

import (
	"log"
	"os"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

func HashObject() {
	sourcePath := utils.GetArgumentsForStage(3)[0]
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		log.Fatal("HashObject could not read file")
	}
	utils.HashObject(content, true)
}
