package main

import (
	"ibse/cmd"
	"ibse/internal"
	"os"
)

func main() {
	internal.Init()

	cmd.Execute()

	os.RemoveAll(internal.TempDir)
}
