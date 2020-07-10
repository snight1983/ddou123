package ddlib

import (
	"sync"

	"github.com/emirpasic/gods/lists/arraylist"
)

//ObjCache .
type ObjCache struct {
	cacheList     *arraylist.List
	cacheListlock sync.Mutex
}

//NewCache .
func NewCache() *ObjCache {
	cache := &ObjCache{}
	cache.cacheList = arraylist.New()
	return cache
}

//GetObj .
func (obj *ObjCache) GetObj() interface{} {
	defer obj.cacheListlock.Unlock()
	obj.cacheListlock.Lock()
	if obj.cacheList.Size() > 0 {
		inter, _ := obj.cacheList.Get(0)
		obj.cacheList.Remove(0)
		return inter
	}
	return nil
}

//BackObj .
func (obj *ObjCache) BackObj(inter interface{}) {
	defer obj.cacheListlock.Unlock()
	obj.cacheListlock.Lock()
	if obj.cacheList.Size() < CACHEMAX {
		obj.cacheList.Add(inter)
	}
}
