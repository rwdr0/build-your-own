package main

import (
	"log"

	"github.com/rwdr0/build-your-own/git/app/packfile"
	"github.com/rwdr0/build-your-own/git/app/utils"
)

// Clone implements the "git clone" command, fetching objects from the given
// remote URL using the git smart HTTP protocol.
func Clone() {
	url := utils.GetArgumentsForStage(2)[0]
	packFile, err := packfile.Fetch(url)
	if err != nil {
		log.Fatalf("Failed to fetch packfile")
	}

	packfile.Unpack(packFile)
}
