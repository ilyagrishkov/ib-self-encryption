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
	Use:   "add [path to file] [data map output dir]",
	Short: "encrypt file and add it to Fabric",
	Long:  `Encrypt file using ID-based self-encryption, upload chunks to IPFS, and create new Fabric asset`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		output, err := internal.Encrypt(args[0])
		if err != nil {
			return
		}
		chunks := getChunkNames(output)
		path, _ := zipChunks(chunks)
		cid, err := internal.SendToIPFS(path)
		if err != nil {
			return
		}
		fabric := internal.NewFabric()
		err = fabric.CreateAsset(generateRandomID(15), cid)
		if err != nil {
			return
		}
		err = internal.CopyFile(fmt.Sprintf("%s/data_map", output), args[1])
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

func getChunkNames(outputPath string) []string {
	var chunks []string
	items, _ := ioutil.ReadDir(outputPath)
	for _, item := range items {
		if item.Name() != "data_map" {
			chunks = append(chunks, fmt.Sprintf("%s/chunk_store/%s", internal.TempDir, item.Name()))
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

func zipChunks(chunks []string) (string, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(fmt.Sprintf("%s/chunks.zip", internal.TempDir), flags, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	zipw := zip.NewWriter(file)
	defer zipw.Close()

	for _, filename := range chunks {
		if err := appendFiles(filename, zipw); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s/chunks.zip", internal.TempDir), nil
}