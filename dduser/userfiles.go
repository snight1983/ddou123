package dduser

import (
	"ddou123/ddlib"
	"strconv"
	"time"

	"github.com/emirpasic/gods/lists/arraylist"
	guuid "github.com/google/uuid"
)

var (
	userFilesCache *ddlib.ObjCache = ddlib.NewCache()
)

// UserFile .
type UserFile struct {
	FileID        string
	UFileID       string
	UFileParentID string
	DisplayName   string
	Describe      string
	IsFolder      bool
	CreateTm      int64
	Year          int
	Month         int
	Day           int
	DelateTm      int64
	VisitNum      int
	LastAccessTm  int64
	Reserve       string
}

func (obj *UserFile) copy(newObj *UserFile) {
	obj.FileID = newObj.FileID
	obj.UFileID = newObj.UFileID
	obj.UFileParentID = newObj.UFileParentID
	obj.DisplayName = newObj.DisplayName
	obj.Describe = newObj.Describe
	obj.IsFolder = newObj.IsFolder
	obj.CreateTm = newObj.CreateTm
	obj.Year = newObj.Year
	obj.Month = newObj.Month
	obj.Day = newObj.Day
	obj.DelateTm = newObj.DelateTm
	obj.VisitNum = newObj.VisitNum
	obj.LastAccessTm = newObj.LastAccessTm
	obj.Reserve = newObj.Reserve
}

func (obj *UserFile) clean() {
	obj.FileID = ""
	obj.UFileID = ""
	obj.UFileParentID = ""
	obj.DisplayName = ""
	obj.Describe = ""
	obj.IsFolder = false
	obj.CreateTm = 0
	obj.Year = 0
	obj.Month = 0
	obj.Day = 0
	obj.DelateTm = 0
	obj.VisitNum = 0
	obj.LastAccessTm = 0
	obj.Reserve = ""
}

// UserFiles .
type UserFiles struct {
	UFile   *UserFile
	UFChild *arraylist.List
}

func initUserfile(uFile *UserFile, fid string, ufid string, parentufid string, name string, desc string, isfolder bool) {
	uFile.FileID = fid
	if len(ufid) == 0 {
		uFile.UFileID = guuid.New().String()
	} else {
		uFile.UFileID = ufid
	}
	uFile.UFileParentID = parentufid
	uFile.DisplayName = name
	uFile.Describe = desc
	uFile.IsFolder = isfolder
	uFile.CreateTm = time.Now().Unix()
	uFile.LastAccessTm = uFile.CreateTm
	uFile.Year, _ = strconv.Atoi(time.Now().Format("2006"))
	uFile.Month, _ = strconv.Atoi(time.Now().Format("01"))
	uFile.Day, _ = strconv.Atoi(time.Now().Format("02"))
	uFile.VisitNum = 0
}

func getUserFiles() *UserFiles {
	ufs := userFilesCache.GetObj()
	if nil == ufs {
		ufs = &UserFiles{}
		ufs.(*UserFiles).UFile = &UserFile{}
		ufs.(*UserFiles).UFChild = arraylist.New()
	}
	return ufs.(*UserFiles)
}

func backUserFiles(ufs *UserFiles) {
	ufs.UFile.clean()
	ufs.UFChild.Clear()
	userFilesCache.BackObj(ufs)
}
