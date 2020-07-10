package dduser

import (
	"ddou123/ddlib"
	"fmt"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"
)

var (
	gUserFilesDBRootPath string
)

func updateuileDB(db *leveldb.DB, ufile *UserFile) error {
	ret, err := msgpack.Marshal(ufile)
	if nil != err {
		ddlib.GLOGER.Error("", zap.String("error", fmt.Sprintf("%s", err)))
		return err
	}
	err = db.Put([]byte(ufile.FileID), []byte(ret), nil)
	if nil != err {
		ddlib.GLOGER.Error("", zap.String("error", fmt.Sprintf("%s", err)))
		return err
	}
	return nil
}

func appendChild(root *UserFiles, name string, fid string, ufid string, desc string, isfolder bool) {
	child := getUserFiles()
	initUserfile(child.UFile, fid, ufid, root.UFile.UFileID, name, desc, true)
	root.UFChild.Add(child)
}

func createRoot(uuid string) *UserFiles {
	ufsRoot := getUserFiles()
	path := filepath.Join(gUserFilesDBRootPath, uuid)
	ufDB, err := leveldb.OpenFile(path, nil)
	if nil != err {
		ddlib.GLOGER.Error("", zap.String("error", fmt.Sprintf("%s", err)))
		return nil
	}
	defer ufDB.Close()
	initUserfile(ufsRoot.UFile, "", ddlib.ROOTFOLDER, "", ddlib.ROOTFOLDER, ddlib.ROOTFOLDER, true)
	updateuileDB(ufDB, ufsRoot.UFile)

	appendChild(ufsRoot, "图片", "", "", "", true)
	appendChild(ufsRoot, "视频", "", "", "", true)
	appendChild(ufsRoot, "音乐", "", "", "", true)
	appendChild(ufsRoot, "文档", "", "", "", true)
	appendChild(ufsRoot, "笔记", "", "", "", true)
	appendChild(ufsRoot, "其它", "", "", "", true)

	return ufsRoot
}
