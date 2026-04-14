package main

import (
	"fmt"

	"github.com/rwdr0/build-your-own/git/app/utils"
)

// CommitTree implements the "git commit-tree" command, creating a commit object
// from a tree SHA, parent SHA, and message, then printing the resulting commit hash.
func CommitTree() {
	args := utils.GetArgumentsForStage(2, 4, 6)
	treeSha := args[0]
	parentSha := args[1]
	commitMessage := args[2]
	author := "John Doe <john@example.com> 1234567890 +0000"
	commiter := "John Doe <john@example.com> 1234567890 +0000"
	commitObject := fmt.Sprintf("tree %s\nparent %s\nauthor %s\ncommiter %s\n\n%s\n", treeSha, parentSha, author, commiter, commitMessage)
	header := fmt.Sprintf("commit %d\x00", len(commitObject))
	commitObject = header + commitObject

	_, hexHash := utils.WriteObject([]byte(commitObject))
	fmt.Println(hexHash)
}
