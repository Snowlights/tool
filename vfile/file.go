package vfile

import (
	"io"
	"os"
	"strings"
)

func FileExists(filename string) bool {

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func IsFile(filename string) (bool, error) {

	if stat, err := os.Stat(filename); err != nil {
		return false, err
	} else {
		return !stat.IsDir(), nil
	}

}

func FileSize(filename string) (int64, error) {

	if stat, err := os.Stat(filename); err != nil {
		return 0, err
	} else if !stat.IsDir() {
		return stat.Size(), nil
	}

	return 0, os.ErrInvalid
}

func LineCount(pathname string) (cnt uint32, err error) {
	var file *os.File
	file, err = os.Open(pathname)
	if err != nil {
		return
	}
	defer file.Close()

	var buf []byte = make([]byte, 1024)
	var rdSize int
	for {
		rdSize, err = file.Read(buf)
		if rdSize > 0 {
			for i := 0; i < rdSize; i++ {
				if buf[i] == '\n' {
					cnt++
				}
			}
		} else {
			if err == io.EOF {
				err = nil
				break
			} else {
				return 0, err
			}
		}
	}
	return
}

func FileCopy(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		destFile.Close()
		return err
	}
	return destFile.Close()
}

func DirExists(p string) bool {

	if file, err := os.Stat(p); err == nil {
		return file.IsDir()
	}
	return false
}

func SearchFileInDir(filename string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if strings.HasSuffix(path, "/") {
			fullpath = path + filename
		} else {
			fullpath = path + "/" + filename
		}
		if FileExists(fullpath) {
			return
		}
	}
	fullpath = ``
	err = os.ErrNotExist
	return
}
