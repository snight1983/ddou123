package dduser

import (
	"ddou123/ddlib"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	gUserBaseDBPath string
	gUserInfoMap    sync.Map        // uuid, userinfo(userbase)
	gUserInfoCache  *ddlib.ObjCache = ddlib.NewCache()
)

//AccountType .
type AccountType int8

const (
	//ACUnKnow .
	ACUnKnow AccountType = 1
	//ACPhone .
	ACPhone AccountType = 2
	//ACEmail .
	ACEmail AccountType = 3
	//ACQQ .
	ACQQ AccountType = 4
	//ACWX .
	ACWX AccountType = 5
)

//LoginStatu .
type LoginStatu int8

const (
	//LoginLive .
	LoginLive LoginStatu = 1
	//LoginOut .
	LoginOut LoginStatu = 2
	//LoginBusy .
	LoginBusy LoginStatu = 3
	//LoginHide .
	LoginHide LoginStatu = 4
)

/*

	UserInfo
	|
	|-------------------UserBase
	|
	|-------------------UserFiles
	|
	|-------------------UserRelationship

*/

//UserInfo .
type UserInfo struct {
	UBase                 *UserBase
	UFiles                sync.Map //(UFileID, *UserFiles)
	UFileRoot             *UserFiles
	UserFileDB            *leveldb.DB
	UserFileDBPath        string
	UserRelationShipDB    *leveldb.DB
	UserRelationShipBPath string
	Token                 string
	IsUserFiles           bool
	IsRelationship        bool
	LoginTime             int64
	LoginOutTime          int64
	LoginAddr             string
	LoginIP               string
	EquipmentID           string
	Version               string
	Longitude             float32
	Lalatitude            float32
	RouteID               string
	LStatu                LoginStatu
}

func (obj *UserInfo) clean() {
	backUserBase(obj.UBase)
	obj.UFiles.Range(func(k, v interface{}) bool {
		backUserFiles(v.(*UserFiles))
		obj.UFiles.Delete(k)
		return true
	})

	obj.UBase = nil
	obj.UFileRoot = nil
	obj.UserFileDB = nil
	obj.UserFileDBPath = ""
	obj.UserRelationShipDB = nil
	obj.UserRelationShipBPath = ""
	obj.Token = ""
	obj.IsUserFiles = false
	obj.IsRelationship = false
	obj.LoginTime = 0
	obj.LoginOutTime = 0
	obj.LStatu = LoginLive
	obj.LoginAddr = ""
	obj.LoginIP = ""
	obj.EquipmentID = ""
	obj.Version = ""
	obj.Longitude = 0.0
	obj.Lalatitude = 0.0
	obj.RouteID = ""
}

func getUserInfo() *UserInfo {
	uinfo := gUserInfoCache.GetObj()
	if nil == uinfo {
		uinfo = &UserInfo{}
	}
	return uinfo.(*UserInfo)
}

func backUserInfo(uinfo *UserInfo) {
	uinfo.clean()
	gUserInfoCache.BackObj(uinfo)
}
