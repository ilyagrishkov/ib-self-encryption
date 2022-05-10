package cmd

import (
	"archive/zip"
	"fmt"
	"github.com/spf13/cobra"
	"ibse/internal"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
)

var addCmd = &cobra.Command{
	Use:   "add [identity] [path to file] [data map output dir]",
	Short: "encrypt file and add it to Fabric",
	Long:  `Encrypt file using ID-based self-encryption, upload chunks to IPFS, and create new Fabric asset`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		output, err := internal.Encrypt(args[1], args[0])
		if err != nil {
			return
		}
		chunks := getChunkNames(output)
		var CIDs []string
		for chunk := range chunks {
			zipped, err := zipChunk(chunk)
			if err != nil {
				return
			}
			cid, _ := internal.SendToIPFS(zipped)
			CIDs = append(CIDs, cid)
		}
		fabric := internal.NewFabric()
		err = fabric.CreateAsset(generateRandomID(15), args[0], CIDs)
		if err != nil {
			return
		}
		err = internal.CopyFile(fmt.Sprintf("%s/data_map", output), args[2])
		if err != nil {
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}

func generateRandomID(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getChunkNames(outputPath string) map[string]bool {
	chunks := map[string]bool{}
	items, _ := ioutil.ReadDir(outputPath)
	for _, item := range items {
		if item.Name() != "data_map" {
			chunks[fmt.Sprintf("%s/chunk_store/%s", internal.TempDir, item.Name())] = true
		}
	}
	return chunks
}

func appendFiles(filename string, zipWriter *zip.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open %s: %s", filename, err)
	}
	defer file.Close()

	wr, err := zipWriter.Create(path.Base(filename))
	if err != nil {
		return fmt.Errorf("failed to create entry for %s in zip file: %s", filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("failed to write %s to zip: %s", filename, err)
	}

	return nil
}

func zipChunk(chunk string) (string, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(fmt.Sprintf("%s/chunks.zip", internal.TempDir), flags, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	zipw := zip.NewWriter(file)
	defer zipw.Close()

	chunkFile, err := os.Open(chunk)
	if err != nil {
		return "", fmt.Errorf("failed to open %s: %s", chunk, err)
	}
	defer chunkFile.Close()

	wr, err := zipw.Create(path.Base(chunk))
	if err != nil {
		return "", fmt.Errorf("failed to create entry for %s in zip file: %s", chunk, err)
	}

	if _, err := io.Copy(wr, chunkFile); err != nil {
		return "", fmt.Errorf("failed to write %s to zip: %s", chunk, err)
	}

	return fmt.Sprintf("%s/chunks.zip", internal.TempDir), nil
}
