package internal

import (
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"io"
	"io/ioutil"
	"os"
	"path"
	"rustgo"
)

func SendToIPFS(filepath string) (string, error) {
	sh := shell.NewShell("localhost:5001")

	file, err := os.Open(filepath)
	if err != nil {
		return "", nil
	}

	return sh.Add(file)
}

func Encrypt(filepath string) (string, error) {
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, TempDir)

	filename := path.Base(filepath)
	dst := fmt.Sprintf("%s/%s", TempDir, filename)

	err := CopyFile(filepath, dst)
	if err != nil {
		return "", err
	}

	_, err = wasm.Invoke("self_encrypt", rustgo.Void, filename)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/chunk_store", TempDir), nil
}

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
