package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"ibse/internal"
)

var readCmd = &cobra.Command{
	Use:   "read [id]",
	Short: "Read asset with given id or all assets",
	Long:  `Encrypt file using ID-based self-encryption, upload chunks to IPFS, and create new Fabric asset`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if all, _ := cmd.Flags().GetBool("all"); len(args) == 1 && all {
			fmt.Println("Provide either id or -a / --all flag")
			return
		}
		fabric := internal.NewFabric()
		if all, _ := cmd.Flags().GetBool("all"); all {
			assets, err := fabric.ReadAllAssets()
			if err != nil {
				return
			}
			fmt.Println(assets)
			return
		}
		if len(args) == 1 {
			asset, err := fabric.ReadAsset(args[0])
			if err != nil {
				fmt.Printf("Error reading asset with id %s\n", args[0])
			}
			fmt.Println(asset)
		} else {
			fmt.Println("Provide asset ID or -a flag to read all assets")
		}
	},
}

func init() {
	RootCmd.AddCommand(readCmd)
	readCmd.Flags().BoolP("all", "a", false, "Read all assets")
}
