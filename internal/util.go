package internal

import (
	"archive/zip"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"rustgo"
)

// SendToIPFS Send data to the IPFS
func SendToIPFS(filepath string) (string, error) {
	sh := shell.NewShell("localhost:5001")

	file, err := os.Open(filepath)
	if err != nil {
		return "", nil
	}

	return sh.Add(file)
}

// GetFromIPFS Receive data from the IPFS
func GetFromIPFS(cid string) (string, error) {
	sh := shell.NewShell("localhost:5001")

	out := fmt.Sprintf("%s/%s.zip", TempDir, cid)
	err := sh.Get(cid, out)
	if err != nil {
		return "", err
	}
	return out, nil
}

// Encrypt Perform ID-based self-encryption on a file
func Encrypt(filepath string, identity string) (string, error) {
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, TempDir)

	filename := path.Base(filepath)
	dst := fmt.Sprintf("%s/%s", TempDir, filename)

	err := CopyFile(filepath, dst)
	if err != nil {
		return "", err
	}

	_, err = wasm.Invoke("self_encrypt", rustgo.Void, filename, identity)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/chunk_store", TempDir), nil
}

// Decrypt call Rust self_encrypt function to decrypt a file and write it to a specified destination
func Decrypt(filepath string, destination string, identity string) (string, error) {
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, TempDir)

	filename := path.Base(destination)
	if filepath != fmt.Sprintf("%s/chunk_store", TempDir) {
		err := CopyDir(filepath, fmt.Sprintf("%s/chunk_store", TempDir))
		if err != nil {
			return "", err
		}
	}

	_, err := wasm.Invoke("self_decrypt", rustgo.Void, filename, identity)
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("%s/%s", TempDir, filename)
	err = CopyFile(output, destination)
	if err != nil {
		return "", err
	}

	return output, nil
}

// CopyFile make a copy of a file
func CopyFile(src string, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// CopyDir Make a copy of a directory
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// Unzip Extract a single file from an archive
func Unzip(src string, dest string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	file := r.File[0]
	filePath := filepath.Join(dest, file.Name)

	if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return "", err
	}

	outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return "", err
	}

	rc, err := file.Open()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(outFile, rc)

	outFile.Close()
	rc.Close()

	if err != nil {
		return "", err
	}

	return filePath, nil
}

// ZipChunk Zip a single chunk
func ZipChunk(chunk string) (string, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(fmt.Sprintf("%s/chunks.zip", TempDir), flags, 0644)
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

	return fmt.Sprintf("%s/chunks.zip", TempDir), nil
}

// GetChunkNames Get encrypted chunks' absolute locations
func GetChunkNames(outputPath string) map[string]bool {
	chunks := map[string]bool{}
	items, _ := ioutil.ReadDir(outputPath)
	for _, item := range items {
		if item.Name() != "data_map" {
			chunks[fmt.Sprintf("%s/chunk_store/%s", TempDir, item.Name())] = true
		}
	}
	return chunks
}

// CreateChunkStore Create a directory chunk_store that contains encrypted chunks and a data map
func CreateChunkStore(dataMapLoc string, chunksLocs []string) {
	err := os.Mkdir(fmt.Sprintf("%s/chunk_store", TempDir), os.ModePerm)
	if err != nil {
		return
	}
	for _, chunkLoc := range chunksLocs {
		_, err := Unzip(chunkLoc, fmt.Sprintf("%s/chunk_store", TempDir))
		if err != nil {
			return
		}
	}

	err = CopyFile(dataMapLoc, fmt.Sprintf("%s/chunk_store/data_map", TempDir))
	if err != nil {
		return
	}
}
