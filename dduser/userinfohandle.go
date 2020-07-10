package dduser

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/vmihailenco/msgpack"
)

func (obj *UserInfo) loadUserRes() {
	if false == obj.IsUserFiles {
		obj.IsUserFiles = true

	}

	if false == obj.IsRelationship {
		obj.IsRelationship = true
	}
}

func (obj *UserInfo) loadUserFile() error {
	var err error
	obj.UserFileDB, err = leveldb.OpenFile(obj.UserFileDBPath, nil)
	defer obj.UserFileDB.Close()
	if nil != err {
		return err
	}

	iter := obj.UserFileDB.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		ufsChild := getUserFiles()
		if err := msgpack.Unmarshal(iter.Value(), ufsChild.UFile); err == nil {
			value, have := obj.UFiles.Load(ufsChild.UFile.UFileID)
			if true == have {
				value.(*UserFiles).UFile.copy(ufsChild.UFile)
				backUserFiles(ufsChild)
				ufsChild = value.(*UserFiles)
			} else {
				obj.UFiles.Store(ufsChild.UFile.UFileID, ufsChild)
			}

			ufsParent, hv := obj.UFiles.Load(ufsChild.UFile.UFileParentID)
			if false == hv {
				ufsParent = getUserFiles()
				ufsParent.(*UserFiles).UFile.FileID = ufsChild.UFile.UFileParentID
				obj.UFiles.Store(ufsParent.(*UserFiles).UFile.FileID, ufsParent)
			}
			ufsParent.(*UserFiles).UFChild.Add(ufsChild.UFile)
			continue
		}
		backUserFiles(ufsChild)
	}
	return iter.Error()
}

func loadUserBase() error {
	userBaseDB, err := leveldb.OpenFile(gUserBaseDBPath, nil)
	defer userBaseDB.Close()
	if nil != err {
		return err
	}
	iter := userBaseDB.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		userBase := getUserBase()
		if err := msgpack.Unmarshal(iter.Value(), userBase); err != nil {
			fmt.Println(err)
			backUserBase(userBase)
		} else {
			userinfo := getUserInfo()
			userinfo.UBase = userBase
			fmt.Println(userinfo.UBase.UUID)
			gUserInfoMap.Store(userBase.UUID, userinfo)
		}
	}

	err = iter.Error()
	if err != nil {
		return err
	}
	return nil
}
