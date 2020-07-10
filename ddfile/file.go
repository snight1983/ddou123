package ddfile

import (
	"ddou123/ddlib"
	"path/filepath"
)

var (
	//CACHEMAX .
	CACHEMAX     int = 8192
	gWorkFolders []string
	gFileFolders []fileFolder
	gDBRootPath  string
)

type fileFolder struct {
	volume string
	path   string
}

//AddFolder 增加文件系统目录
func AddFolder(path string) {
	vname := filepath.VolumeName(path)
	for _, value := range gFileFolders {
		if value.volume == vname {
			return
		}
	}
	if nil != ddlib.IsExist(path) {
		if nil != ddlib.CreateMutiDir(path) {
			return
		}
	}
	gFileFolders = append(gFileFolders, fileFolder{vname, path})
}
