package main

import (
	"github.com/rwdr0/build-your-own/git/app/packfile"
	"github.com/rwdr0/build-your-own/git/app/utils"
)

func Clone() {
	url := utils.GetArgumentsForStage(2)[0]
	packfile.Fetch(url)
}
