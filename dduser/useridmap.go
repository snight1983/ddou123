package dduser

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	gUserIDMap       sync.Map // key, uuid
	gUserIDMapDBPath string
)

func initUserIDMap() error {
	userIDDB, err := leveldb.OpenFile(gUserIDMapDBPath, nil)
	defer userIDDB.Close()
	if nil != err {
		return err
	}

	iter := userIDDB.NewIterator(nil, nil)
	defer iter.Release()
	for iter.Next() {
		gUserIDMap.Store(string(iter.Key()), string(iter.Value()))
	}
	err = iter.Error()
	if err != nil {
		return err
	}
	return nil
}
