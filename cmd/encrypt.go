package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"ibse/internal"
	"path"
	"rustgo"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt [path to file] [output path]",
	Short: "encrypts file locally",
	Long:  `This subcommand encrypts file using ID-based self-encryption`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := encrypt(args[0], args[1])
		if err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(encryptCmd)
}

func encrypt(filepath string, destination string) error {
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, internal.TempDir)

	filename := path.Base(filepath)
	dst := fmt.Sprintf("%s/%s", internal.TempDir, filename)

	err := internal.CopyFile(filepath, dst)
	if err != nil {
		return err
	}

	_, err = wasm.Invoke("self_encrypt", rustgo.Void, filename)
	if err != nil {
		return err
	}

	output := fmt.Sprintf("%s/chunk_store", internal.TempDir)
	err = internal.CopyDir(output, destination)
	if err != nil {
		return err
	}

	return nil
}
