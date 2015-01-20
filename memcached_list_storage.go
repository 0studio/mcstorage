package storage

import (
	"github.com/0studio/storage_key"
)

var MAXLEN = 200

func (this MemcachedStorage) Getlimit(key key.Key, sinceId, maxId interface{}, page, count int) (interface{}, error) {
	obj, err := this.Get(key)
	if err != nil {
		return nil, err
	}
	return Page(obj.(Pagerable), sinceId, maxId, page, count), nil
}

func (this MemcachedStorage) AddItem(key key.Key, item interface{}) error {
	obj, err := this.Get(key)
	if err != nil {
		return err
	}
	result := obj.(Pagerable).AddItem(item, MAXLEN)
	return this.Set(key, result)
}

func (this MemcachedStorage) DeleteItem(key key.Key, item interface{}) error {
	obj, err := this.Get(key)
	if err != nil {
		return err
	}
	result := obj.(Pagerable).DeleteItem(item)
	return this.Set(key, result)
}
