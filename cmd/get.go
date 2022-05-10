package cmd

import (
	"fmt"
	"github.com/cloudflare/cfssl/log"
	"github.com/spf13/cobra"
	"ibse/internal"
)

var getCmd = &cobra.Command{
	Use:   "get [block id] [path to data map] [destination path]",
	Short: "get file from Fabric and decrypt it",
	Long:  `Get chunks of encrypted file from IPFS, and decrypt them`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		// Create new Fabric connection
		fabric := internal.NewFabric()

		// Read asset from the Fabric
		asset, err := fabric.ReadAsset(args[0])
		if err != nil {
			log.Error("Failed to read an asset from the Fabric")
			return
		}

		// Iterate over CIDs and download all file chunks from the IPFS
		var paths []string
		for _, cid := range asset.CID {
			path, _ := internal.GetFromIPFS(cid)
			paths = append(paths, path)
		}

		// Generate a chunk_store directory in a WASM mapped directory
		internal.CreateChunkStore(args[1], paths)

		// Decrypt a file and write the output to the specified location
		_, err = internal.Decrypt(fmt.Sprintf("%s/chunk_store", internal.TempDir), args[2], fabric.PublicKey)
		if err != nil {
			log.Error("failed to decrypt a file")
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
