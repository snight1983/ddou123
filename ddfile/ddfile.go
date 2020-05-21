package ddfile

import (
	"bytes"
	"crypto/sha1"
	"ddou123/crt"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/exp/mmap"
)

/*
0-root
	|----1-tmp
	|----1-dime

*/

var gFilesMap sync.Map
var gFolders []string

// Init ...
func Init(paths []string) error {
	for _, path := range paths {
		crt.CreateMutiDir(path)
		if crt.IsExist(path) {
			gFolders = append(gFolders, path)
			scannFiles(path, 0)
		}
	}
	return nil
}

// GetCacheFolder .
func GetCacheFolder() (string, int64) {
	var sizeMax int64 = 0
	var ph string = ""
	for _, path := range gFolders {
		space, err := crt.GetFreeSpace(path)
		if nil == err {
			if sizeMax < space {
				sizeMax = space
				ph = path
			}
		}
	}
	ph = filepath.Join(ph, "tmp")
	crt.CreateMutiDir(ph)
	return ph, sizeMax
}

func getBestSavePath(size int, name string) string {
	var ph string = ""
	var sizeMax int64 = 0

	for _, path := range gFolders {
		space, err := crt.GetFreeSpace(path)
		if nil == err {
			if space > int64(size) {
				if sizeMax < space {
					sizeMax = space
					ph = path
				}
			}
		}
	}

	dailyFolder := time.Now().Format("2006-01-02")
	dailyFolder = filepath.Join(ph, dailyFolder)
	crt.CreateMutiDir(dailyFolder)
	return filepath.Join(dailyFolder, name)
}

// GetPath .
func GetPath(key string) (string, bool) {

	value, ishave := gFilesMap.Load(key)
	if ishave {
		path := value.(string)
		return path, crt.IsExist(path)
	}
	return "", ishave
}

func getFileSha(pos int64, total int64, at *mmap.ReaderAt) (string, bool) {
	buff := make([]byte, 1024*1024*8)
	var bufSha bytes.Buffer
	for {
		lenRead, err := at.ReadAt(buff, pos)
		pos += int64(lenRead)
		h := sha1.New()
		h.Write(buff)
		bufSha.Write(h.Sum(nil))
		if pos == total {
			bufSha.Bytes()
			h := sha1.New()
			h.Write(buff)
			return fmt.Sprintf("%x", h.Sum(nil)), true
		}
		if nil != err {
			return "", false
		}
	}
}

// SaveFile .
func SaveFile(path string) (string, bool) {
	if crt.IsExist(path) {
		file, err := mmap.Open(path)
		if nil != err {
			return "", false
		}
		defer file.Close()

		sha, res := getFileSha(int64(0), int64(file.Len()), file)
		if false == res {
			return "", false
		}

		if ph, ok := GetPath(sha); true == ok {
			return ph, ok
		}
		ph := getBestSavePath(file.Len(), sha)
		crt.CopyFile(int64(0), int64(file.Len()), ph, file)
		res = crt.IsExist(ph)
		if res {
			gFilesMap.Store(sha, ph)
		}
		return ph, res
	}
	return "", false
}

func scannFiles(path string, depth int) {
	fs, _ := ioutil.ReadDir(path)
	for _, file := range fs {
		if file.IsDir() {
			if depth == 0 {
				fullName := file.Name()
				if _, err := time.Parse("2006-01-02", fullName); err == nil {
					path := filepath.Join(path, fullName)
					fmt.Println("folder:" + path)
					scannFiles(path, (depth + 1))
				}
			}
		} else {
			if depth == 1 {
				gFilesMap.Store(file.Name(), filepath.Join(path, file.Name()))
				fmt.Println("key:" + file.Name())
				fmt.Println("value:" + filepath.Join(path, file.Name()))
			}
		}
	}
}
