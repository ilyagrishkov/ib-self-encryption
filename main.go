package main

import (
	"ibse/cmd"
	"ibse/internal"
	"os"
)

func main() {
	internal.Init()

	cmd.Execute()

	cleanUp()
}

func cleanUp() {
	_ = os.RemoveAll(internal.TempDir)
}
