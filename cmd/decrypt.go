package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"ibse/internal"
	"path"
	"rustgo"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt [identity] [path to directory with data map and file chunks] [output path]",
	Short: "decrypts file locally",
	Long:  `This subcommand decrypts ID-based self-encrypted file`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		err := decrypt(args[1], args[2], args[0])
		if err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(decryptCmd)
}

func decrypt(filepath string, destination string, identity string) error {
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", internal.RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, internal.TempDir)

	filename := path.Base(destination)
	err := internal.CopyDir(filepath, fmt.Sprintf("%s/chunk_store", internal.TempDir))
	if err != nil {
		return err
	}

	_, err = wasm.Invoke("self_decrypt", rustgo.Void, filename, identity)
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
