package ddfile

import (
	"bytes"
	"crypto/sha1"
	"ddou123/ddlib"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"
	"golang.org/x/exp/mmap"
)

var (
	//BLOCKSIZE .
	BLOCKSIZE       int64 = 4194304 // 1024 * 1024 * 4
	gFileBaseMap    sync.Map
	gFileBaseDB     *leveldb.DB
	gFileBaseDBPath string
	fileBaseCache   *ddlib.Queue = ddlib.NewQueue()
)

// FileBase .
type FileBase struct {
	FileID     string
	Path       string
	Ext        string
	Size       int64
	CrTime     int64
	RefCnt     int32
	LastAcTime int64
}

func (obj *FileBase) clean() {
	obj.FileID = ""
	obj.Path = ""
	obj.Ext = ""
	obj.Size = 0
	obj.CrTime = 0
	obj.RefCnt = 0
}

func getFileBase() *FileBase {
	ub := fileBaseCache.Pop()
	if nil == ub {
		ub = &FileBase{}
	}
	return ub.(*FileBase)
}

func backFileBase(obj *FileBase) {
	obj.clean()
	fileBaseCache.Push(obj)
}

func initFB() error {
	var err error = nil
	gFileBaseDB, err = leveldb.OpenFile(gFileBaseDBPath, nil)
	if nil != err {
		ddlib.GLoger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
		return err
	}
	delList := arraylist.New()
	iter := gFileBaseDB.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		fileBase := getFileBase()
		if err := msgpack.Unmarshal(iter.Value(), fileBase); err != nil {
			ddlib.GLoger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
		} else {
			if err := ddlib.IsExist(fileBase.Path); nil == err {
				if fileBase.RefCnt == 0 {
					delList.Add(fileBase)
					continue
				}
				backFileBase(fileBase)
				continue
			}
		}
		delList.Add(fileBase)
	}

	err = iter.Error()
	if err != nil {
		ddlib.GLoger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
		return err
	}

	it := delList.Iterator()
	for it.Next() {
		fb := it.Value().(*FileBase)
		gFileBaseDB.Delete([]byte(fb.FileID), nil)
		delateFileByPath(fb.Path)
		backFileBase(fb)
	}
	return nil
}

func uninitFB() error {
	return gFileBaseDB.Close()
}

func updateFB(obj *FileBase) error {
	ret, err := msgpack.Marshal(obj)
	if err == nil {
		return gFileBaseDB.Put([]byte(obj.FileID), []byte(ret), nil)
	}
	return err
}

func delateFB(fb *FileBase) error {
	fb.RefCnt--
	if fb.RefCnt > 0 {
		return updateFB(fb)
	}
	delateFileByPath(fb.Path)
	err := gFileBaseDB.Delete([]byte(fb.FileID), nil)
	gFileBaseMap.Delete(fb.FileID)
	backFileBase(fb)
	return err
}

func getFB(fileid string) (*FileBase, error) {
	value, ishave := gFileBaseMap.Load(fileid)
	if ishave {
		value.(*FileBase).LastAcTime = time.Now().Unix()
		return value.(*FileBase), ddlib.IsExist(value.(*FileBase).Path)
	}

	buff, err := gFileBaseDB.Get([]byte(fileid), nil)
	if nil != err {
		ddlib.GLoger.Warn("", zap.String("error", fmt.Sprintf("%s", err)))
		return nil, err
	}

	fb := getFileBase()
	if err = msgpack.Unmarshal(buff, fb); err != nil {
		goto onExit
	}

	if err = ddlib.IsExist(fb.Path); nil != err {
		backFileBase(fb)
		goto onExit
	}

	fb.LastAcTime = time.Now().Unix()
	gFileBaseMap.Store(fb.FileID, fb)
	return fb, nil
onExit:
	ddlib.GLoger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
	return nil, err
}

// GetTempFolder .
func GetTempFolder() (string, int64) {
	var sizeMax int64 = 0
	var ph string = ""
	for _, path := range gWorkFolders {
		space, err := ddlib.GetFreeSpace(path)
		if nil == err {
			if sizeMax < space {
				sizeMax = space
				ph = path
			}
		}
	}
	ph = filepath.Join(ph, "temp")
	ddlib.CreateMutiDir(ph)
	return ph, sizeMax
}

func getBestSaveFolder(size int64) string {
	var ph string = ""
	var sizeMax int64 = 0

	for _, path := range gWorkFolders {
		space, err := ddlib.GetFreeSpace(path)
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
	ddlib.CreateMutiDir(dailyFolder)
	return dailyFolder
}

func getFileIDBySha(pos int64, total int64, at *mmap.ReaderAt) (string, error) {
	buff := make([]byte, BLOCKSIZE)
	var bufSha bytes.Buffer
	for {
		lenRead, err := at.ReadAt(buff, pos)
		pos += int64(lenRead)
		if nil != err && pos != total {
			ddlib.GLoger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
			return "", nil
		}

		h := sha1.New()
		h.Write(buff[:lenRead])
		bufSha.Write(h.Sum(nil))

		if pos == total {
			h := sha1.New()
			h.Write(bufSha.Bytes())
			return fmt.Sprintf("%x", h.Sum(nil)), nil
		}
	}
}

// SaveFile .
func SaveFile(filePath string) (*FileBase, error) {

	if err := ddlib.IsExist(filePath); nil != err {
		ddlib.GLoger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
		return nil, err
	}

	file, err := mmap.Open(filePath)
	if nil != err {
		ddlib.GLoger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
		return nil, err
	}

	defer file.Close()

	fb := getFileBase()
	fb.CrTime = time.Now().Unix()
	fb.Ext = filepath.Ext(filePath)

	fileID, err := getFileIDBySha(int64(0), int64(file.Len()), file)
	if nil != err {
		backFileBase(fb)
		ddlib.GLoger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
		return nil, err
	}

	if ofb, _ := getFB(fileID); nil != ofb {
		backFileBase(fb)
		return ofb, nil
	}

	path := filepath.Join(getBestSaveFolder(int64(file.Len())), fileID)
	err = ddlib.CreateMutiDir(path)
	if nil != err {
		ddlib.GLoger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
		return nil, err
	}

	fnWithSuffix := filepath.Base(filePath)
	path = filepath.Join(path, fnWithSuffix)
	ddlib.CopyFile(int64(0), int64(file.Len()), path, file)

	if err = ddlib.IsExist(path); nil != err {
		backFileBase(fb)
		ddlib.GLoger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
		return nil, err
	}
	fb.FileID = fileID
	fb.Path = path
	fb.Size = int64(file.Len())

	err = updateFB(fb)
	if nil != err {
		backFileBase(fb)
		ddlib.GLoger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
		return nil, err
	}
	gFileBaseMap.Store(fileID, fb)
	return fb, nil
}

// SaveFileFolder .
func SaveFileFolder(folderPath string, cnt int) (string, error) {

	err := ddlib.IsExist(folderPath)
	if nil == err {
		size := int64(cnt) * BLOCKSIZE
		fd := getBestSaveFolder(size)

		h := sha1.New()
		h.Write([]byte(folderPath))
		tmpN := fmt.Sprintf("%x", h.Sum(nil))
		tmpPN := filepath.Join(fd, tmpN)
		tmpPN += ".ftmp"

		destTmp, err := os.Create(tmpPN)
		if nil != err {
			return "", err
		}

		var bufSha bytes.Buffer
		var buf []byte
		var info os.FileInfo

		for i := 0; i < cnt; i++ {
			file, err := os.Open(filepath.Join(folderPath, strconv.Itoa(i)))
			if err != nil {
				destTmp.Close()
				return "", err
			}

			defer file.Close()

			if info, err = file.Stat(); nil != err {
				destTmp.Close()
				return "", err
			}

			len := info.Size()
			if len <= 0 || len > BLOCKSIZE {
				return "", fmt.Errorf("ssss")
			}

			if buf, err = ioutil.ReadAll(file); err != nil {
				destTmp.Close()
				return "", err
			}

			if _, err := destTmp.Write(buf[:len]); err != nil {
				destTmp.Close()
				return "", err
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

		if err = os.Rename(tmpPN, destPN); nil != err {
			return "", err
		}
		if err = ddlib.IsExist(destPN); nil != err {
			return "", err
		}

		gFileBaseMap.Store(sha, destPN)
		return destPN, nil
	}
	return "", err
}

/*
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
				gFileBaseMap.Store(file.Name(), filepath.Join(path, file.Name()))
				fmt.Println("key:" + file.Name())
				fmt.Println("value:" + filepath.Join(path, file.Name()))
			}
		}
	}
}
*/

func delateFileByPath(path string) {
	err := os.Remove(path)
	if nil != err {
		ddlib.GLoger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
	}
}
