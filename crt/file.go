package crt

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/exp/mmap"
)

// IsExist  .
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//CreateMutiDir .
func CreateMutiDir(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败,error info:", err)
			return err
		}
		return err
	}
	return nil
}

// CopyFile .
func CopyFile(pos int64, total int64, dest string, at *mmap.ReaderAt) error {

	var tmp bytes.Buffer
	tmp.WriteString(dest)
	tmp.WriteString(".tmp")
	destTmp, err := os.Create(tmp.String())
	if nil != err {
		return err
	}

	buff := make([]byte, 1024*1024*8)
	for {
		lenRead, err := at.ReadAt(buff, pos)
		pos += int64(lenRead)

		if _, err := destTmp.Write(buff[:lenRead]); err != nil {
			destTmp.Close()
			return err
		}
		if pos == total {
			destTmp.Close()
			return os.Rename(tmp.String(), dest)
		}
		if nil != err {
			destTmp.Close()
			return err
		}
	}
}
