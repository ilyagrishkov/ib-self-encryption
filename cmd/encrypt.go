package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path"
	"rustgo"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt [path to file]",
	Short: "encrypts file locally",
	Long:  `This subcommand encrypts file using ID-based self-encryption`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := encrypt(args[0])
		if err != nil {
			return
		}
		getBool, err := cmd.Flags().GetBool("compress")
		if err != nil {
			return
		}
		fmt.Println(getBool)
	},
}

func init() {
	RootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().BoolP("compress", "c", false, "Compress result")
}

func encrypt(filepath string) error {
	wasm := rustgo.NewWasmLib("rustgo/ib_self_encryption_rust.wasm")
	sourceFileStat, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", filepath)
	}

	source, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer source.Close()

	dst := path.Base(filepath)
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}
	_, err = wasm.Invoke("self_encrypt", rustgo.Void, dst)

	err = os.Remove(dst)
	if err != nil {
		return err
	}

	return nil
}
