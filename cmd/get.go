package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"ibse/internal"
)

var getCmd = &cobra.Command{
	Use:   "get [identity] [block id] [path to data map] [destination path]",
	Short: "get file from Fabric and decrypt it",
	Long:  `Get chunks of encrypted file from IPFS, and decrypt them`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		fabric := internal.NewFabric()
		asset, err := fabric.ReadAsset(args[1])
		if err != nil {
			return
		}

		cid := asset["CID"].(string)
		path, err := internal.GetFromIPFS(cid)
		if err != nil {
			return
		}
		createChunkStore(args[2], path)
		_, err = internal.Decrypt(fmt.Sprintf("%s/chunk_store", internal.TempDir), args[3], args[0])
		if err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}

func createChunkStore(dataMapLoc string, chunksLoc string) {
	_, err := internal.Unzip(chunksLoc, fmt.Sprintf("%s/chunk_store", internal.TempDir))
	if err != nil {
		return
	}
	err = internal.CopyFile(dataMapLoc, fmt.Sprintf("%s/chunk_store/data_map", internal.TempDir))
	if err != nil {
		return
	}
}
