package internal

import (
	"fmt"
	"os"
)

var RootDir string
var TempDir string
var ConfigFile string

func Init() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	RootDir = fmt.Sprintf("%s/.ibse", home)
	TempDir = fmt.Sprintf("%s/.ibse/tmp", home)
	ConfigFile = fmt.Sprintf("%s/.ibse/config.yaml")

	if _, err := os.Stat(TempDir); os.IsNotExist(err) {
		err := os.MkdirAll(TempDir, os.ModePerm)
		if err != nil {
			return
		}
	}

	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		_, err := os.Create(ConfigFile)
		if err != nil {
			return
		}
	}
}
