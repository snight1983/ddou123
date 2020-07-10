package dduser

import (
	"ddou123/ddlib"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/vmihailenco/msgpack"
)

var userBaseCache *ddlib.ObjCache = ddlib.NewCache()

// UserBase 用户基本信息
type UserBase struct {
	UUID           string
	PASSWD         string
	HeaderImg      string
	Name           string
	Sign           string
	Sex            int8
	Age            int8
	Email          string
	Phone          string
	QQ             string
	WX             string
	Address        string
	OnLineDuration int64
	LoginTimes     int32
	ChannelID      string
	Actype         AccountType
}

func (obj *UserBase) clean() {
	obj.UUID = ""
	obj.PASSWD = ""
	obj.HeaderImg = ""
	obj.Name = ""
	obj.Sign = ""
	obj.Sex = 0
	obj.Age = 0
	obj.Email = ""
	obj.Phone = ""
	obj.QQ = ""
	obj.WX = ""
	obj.Address = ""
	obj.Actype = 0
	obj.LoginTimes = 0
	obj.ChannelID = ""
}

func getUserBase() *UserBase {
	ub := userBaseCache.GetObj()
	if nil == ub {
		ub = &UserBase{}
	}
	return ub.(*UserBase)
}

func backUserBase(obj *UserBase) {
	obj.clean()
	userBaseCache.BackObj(obj)
}

func updateUserBase(db *leveldb.DB, obj *UserBase) error {
	ret, err := msgpack.Marshal(obj)
	if err == nil {
		return db.Put([]byte(obj.UUID), []byte(ret), nil)
	}
	return err
}

func delUserBase(db *leveldb.DB, obj *UserBase) error {
	return db.Delete([]byte(obj.UUID), nil)
}
