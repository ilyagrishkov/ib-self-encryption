package cmd

import (
	"github.com/spf13/cobra"
	"ibse/internal"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt [identity] [path to directory with data map and file chunks] [output path]",
	Short: "decrypts file locally",
	Long:  `This subcommand decrypts ID-based self-encrypted file`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := internal.Decrypt(args[1], args[2], args[0])
		if err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(decryptCmd)
}
