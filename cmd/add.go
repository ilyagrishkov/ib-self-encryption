package cmd

import (
	"fmt"
	"github.com/cloudflare/cfssl/log"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"ibse/internal"
)

var addCmd = &cobra.Command{
	Use:   "add [path to file] [data map output dir]",
	Short: "encrypt a file and add it to the Fabric",
	Long:  `Encrypt a file using ID-based self-encryption, upload chunks to IPFS, and create a new Fabric asset`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Create new Fabric connection
		fabric := internal.NewFabric()

		// Encrypt file using ID-based self-encryption with Public Key as an identity
		output, err := internal.Encrypt(args[0], fabric.PublicKey)
		if err != nil {
			return
		}

		// Get encrypted chunk locations
		chunks := internal.GetChunkNames(output)

		// Send each chunk to the IPFS and retrieve their CIDs
		var CIDs []string
		for chunk := range chunks {
			zipped, err := internal.ZipChunk(chunk)
			if err != nil {
				return
			}
			cid, _ := internal.SendToIPFS(zipped)
			CIDs = append(CIDs, cid)
		}

		// Create and asset in Fabric
		err = fabric.CreateAsset(uuid.New().String(), CIDs)
		if err != nil {
			log.Error("Failed to create an asset")
			return
		}

		// Copy data map to the specified location
		err = internal.CopyFile(fmt.Sprintf("%s/data_map", output), args[1])
		if err != nil {
			log.Error("Failed to copy data map to specified location")
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
