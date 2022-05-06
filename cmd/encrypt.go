package cmd

import (
	"github.com/spf13/cobra"
	"ibse/internal"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt [path to file] [output path]",
	Short: "encrypts file locally",
	Long:  `This subcommand encrypts file using ID-based self-encryption`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		output, err := internal.Encrypt(args[0])
		if err != nil {
			return
		}
		err = internal.CopyDir(output, args[1])
		if err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(encryptCmd)
}
