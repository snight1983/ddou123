package ddcloud

import (
	"ddou123/ddlib"
	"errors"
	"path/filepath"

	"go.uber.org/zap/zapcore"
)

var (
	//CACHEMAX .
	CACHEMAX int = 8192

	gWorkFolders []string
	gDBRootPath  string
)

// Init ...
func Init(paths []string) error {

	var err error = nil
	if len(paths) <= 0 {
		return errors.New("设置工作目录")
	}

	if nil != ddlib.CreateMutiDir(paths[0]) {
		return errors.New("工作目录无效")
	}

	if nil == ddlib.GLoger {
		ph := filepath.Join(paths[0], "ddcloud.log")
		ddlib.GLoger = ddlib.NewLogger(ph, zapcore.InfoLevel, 128, 30, 7, true, "ddcloud")
	}

	gDBRootPath = filepath.Join(paths[0], "db")
	if nil != ddlib.CreateMutiDir(gDBRootPath) {
		return errors.New("数据库目录无效")
	}
	return err
	/*
				gFileBaseDBPath = filepath.Join(gDBRootPath, "filebasedb")
				if err = initFB(); err != nil {
					return err
				}

					gUserIDMapDBPath = filepath.Join(gDBRootPath, "useridmapdb")
					if err = initUserIDMap(); err != nil {
						return err
					}


						gUserBaseDBPath = filepath.Join(gDBRootPath, "userbasedb")
						if err = loadUserBase(); err != nil {
							Glogger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
							return err
						}

						gUserFilesDBRootPath = filepath.Join(gDBRootPath, "userfilesdbroot")
						if err = ddlib.CreateMutiDir(gUserFilesDBRootPath); err != nil {
							Glogger.Error("", zap.String("error", fmt.Sprintf("%s", err)))
							return err
						}

			for _, path := range paths {
				ddlib.CreateMutiDir(path)
				if nil == ddlib.IsExist(path) {
					gWorkFolders = append(gWorkFolders, path)
					//scannFiles(path, 0)
				}
			}

		return nil
	*/
}
