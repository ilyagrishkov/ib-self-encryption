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
	"strings"
)

func SendToIPFS(filepath string) (string, error) {
	sh := shell.NewShell("localhost:5001")

	file, err := os.Open(filepath)
	if err != nil {
		return "", nil
	}

	return sh.Add(file)
}

func GetFromIPFS(cid string) (string, error) {
	sh := shell.NewShell("localhost:5001")

	out := fmt.Sprintf("%s/out.zip", TempDir)
	err := sh.Get(cid, out)
	if err != nil {
		return "", err
	}
	return out, nil
}

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

func Decrypt(filepath string, destination string) (string, error) {
	wasmLoc := fmt.Sprintf("%s/ib_self_encryption_rust.wasm", RootDir)
	wasm := rustgo.NewWasmLib(wasmLoc, TempDir)

	filename := path.Base(destination)
	if filepath != fmt.Sprintf("%s/chunk_store", TempDir) {
		err := CopyDir(filepath, fmt.Sprintf("%s/chunk_store", TempDir))
		if err != nil {
			return "", err
		}
	}

	_, err := wasm.Invoke("self_decrypt", rustgo.Void, filename)
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

func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
