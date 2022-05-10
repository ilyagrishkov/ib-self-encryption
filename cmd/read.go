package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/cloudflare/cfssl/log"
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
				log.Error("Failed to read all assets from the Fabric")
				return
			}

			jsonAssets, _ := json.Marshal(assets)
			fmt.Println(string(jsonAssets))
			return
		}
		if len(args) == 1 {
			asset, err := fabric.ReadAsset(args[0])
			if err != nil {
				fmt.Printf("Error reading asset with id %s\n", args[0])
			}
			jsonAsset, _ := json.Marshal(asset)
			fmt.Println(string(jsonAsset))
		} else {
			fmt.Println("Provide asset ID or -a flag to read all assets")
		}
	},
}

func init() {
	RootCmd.AddCommand(readCmd)
	readCmd.Flags().BoolP("all", "a", false, "Read all assets")
}
