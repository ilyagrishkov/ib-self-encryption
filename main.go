package main

import (
	"fmt"
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
	_ = os.RemoveAll(fmt.Sprintf("%s/wallet", internal.RootDir))
}
