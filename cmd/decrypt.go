package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"ibse/internal"
	"path"
	"rustgo"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt [path to directory with data map and file chunks] [output path]",
	Short: "decrypts file locally",
	Long:  `This subcommand decrypts ID-based self-encrypted file`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := decrypt(args[0], args[1])
		if err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(decryptCmd)
}

func decrypt(filepath string, destination string) error {
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, internal.TempDir)

	filename := path.Base(destination)
	err := internal.CopyDir(filepath, fmt.Sprintf("%s/chunk_store", internal.TempDir))
	if err != nil {
		return err
	}

	_, err = wasm.Invoke("self_decrypt", rustgo.Void, filename)
	if err != nil {
		return err
	}

	output := fmt.Sprintf("%s/%s", internal.TempDir, filename)
	err = internal.CopyFile(output, destination)
	if err != nil {
		return err
	}

	return nil
}
