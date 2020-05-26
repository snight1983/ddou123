package ddfile

import (
	"bytes"
	"crypto/sha1"
	"ddou123/crt"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"golang.org/x/exp/mmap"
)

/*
0-root
	|----1-tmp
	|----1-dime

*/

var (
	gFilesMap sync.Map
	gFolders  []string
)

//BLOCKSIZE .
var BLOCKSIZE int64 = 4194304 // 1024 * 1024 * 4

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

func getBestSaveFolder(size int64) string {
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
	return dailyFolder
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
	buff := make([]byte, BLOCKSIZE)
	var bufSha bytes.Buffer
	for {
		lenRead, err := at.ReadAt(buff, pos)
		pos += int64(lenRead)

		h := sha1.New()
		h.Write(buff[:lenRead])
		bufSha.Write(h.Sum(nil))

		if pos == total {
			h := sha1.New()
			h.Write(bufSha.Bytes())
			return fmt.Sprintf("%x", h.Sum(nil)), true
		}
		if nil != err {
			return "", false
		}
	}
}

// SaveFileSingle .
func SaveFileSingle(filePath string) (string, bool) {
	if crt.IsExist(filePath) {
		file, err := mmap.Open(filePath)
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
		ph := filepath.Join(getBestSaveFolder(int64(file.Len())), sha)
		crt.CopyFile(int64(0), int64(file.Len()), ph, file)
		if crt.IsExist(ph) {
			gFilesMap.Store(sha, ph)
			return ph, res
		}
	}
	return "", false
}

// SaveFileFolder .
func SaveFileFolder(folderPath string, cnt int) (string, bool) {

	if crt.IsExist(folderPath) {
		size := int64(cnt) * BLOCKSIZE
		fd := getBestSaveFolder(size)

		h := sha1.New()
		h.Write([]byte(folderPath))
		tmpN := fmt.Sprintf("%x", h.Sum(nil))
		tmpPN := filepath.Join(fd, tmpN)
		tmpPN += ".ftmp"

		destTmp, err := os.Create(tmpPN)
		if nil != err {
			return "", false
		}

		var bufSha bytes.Buffer
		var buf []byte
		var info os.FileInfo

		for i := 0; i < cnt; i++ {
			file, err := os.Open(filepath.Join(folderPath, strconv.Itoa(i)))
			if err != nil {
				destTmp.Close()
				return "", false
			}

			defer file.Close()

			if info, err = file.Stat(); nil != err {
				destTmp.Close()
				return "", false
			}

			len := info.Size()
			if len <= 0 || len > BLOCKSIZE {
				return "", false
			}

			if buf, err = ioutil.ReadAll(file); err != nil {
				destTmp.Close()
				return "", false
			}

			if _, err := destTmp.Write(buf[:len]); err != nil {
				destTmp.Close()
				return "", false
			}

			h := sha1.New()
			h.Write(buf[:len])
			bufSha.Write(h.Sum(nil))
		}
		destTmp.Close()

		h = sha1.New()
		h.Write(bufSha.Bytes())
		sha := fmt.Sprintf("%x", h.Sum(nil))
		destPN := filepath.Join(fd, sha)

		if err := os.Rename(tmpPN, destPN); nil != err {
			return "", false
		}

		if crt.IsExist(destPN) {
			gFilesMap.Store(sha, destPN)
			return destPN, true
		}
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
