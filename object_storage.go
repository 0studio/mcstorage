package storage

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/0studio/storage_key"
)

type Storage interface {
	Get(key key.Key) (interface{}, error)
	Set(key key.Key, object interface{}) error
	Add(key key.Key, object interface{}) error
	MultiGet(keys []key.Key) (map[key.Key]interface{}, error)
	MultiSet(map[key.Key]interface{}) error
	Delete(key key.Key) error
	FlushAll()
}

type CounterStorage interface {
	Storage
	Incr(key key.Key, step uint64) (newValue uint64, err error)
	Decr(key key.Key, step uint64) (newValue uint64, err error)
}

type ListStorage interface {
	Storage
	Getlimit(key key.Key, sinceId, maxId interface{}, page, count int) (interface{}, error)
	AddItem(key key.Key, item interface{}) error
	DeleteItem(key key.Key, item interface{}) error
}

type StorageProxy struct {
	PreferedStorage Storage
	BackupStorage   Storage
}

type Encoding interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte) (interface{}, error)
}

type JsonEncoding struct {
	T reflect.Type
}

func (this JsonEncoding) Marshal(v interface{}) ([]byte, error) {
	buf, err := json.Marshal(v)
	return buf, err
}

func (this JsonEncoding) Unmarshal(data []byte) (interface{}, error) {
	tStruct := reflect.New(this.T)
	dec := json.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(tStruct.Interface())
	if err != nil {
		return err, nil
	}
	return reflect.Indirect(tStruct.Elem()).Interface(), nil
}

type GobEncoding struct {
	T reflect.Type
}

func (this GobEncoding) Marshal(v interface{}) (data []byte, err error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err = enc.Encode(v)
	data = network.Bytes()
	return
}

func (this GobEncoding) Unmarshal(data []byte) (v interface{}, err error) {
	var network bytes.Buffer
	tStruct := reflect.New(this.T)
	network.Write(data)
	dec := gob.NewDecoder(&network)
	err = dec.Decode(tStruct.Interface())
	if err != nil {
		fmt.Println(err)
		return
	}

	v = reflect.Indirect(tStruct.Elem()).Interface()
	return
}

type ByteEncoding struct {
}

func (this ByteEncoding) Marshal(v interface{}) (data []byte, err error) {
	data = v.([]byte)
	return
}

func (this ByteEncoding) Unmarshal(data []byte) (v interface{}, err error) {
	v = data
	return
}

func NewStorageProxy(prefered, backup Storage) StorageProxy {
	return StorageProxy{
		PreferedStorage: prefered,
		BackupStorage:   backup,
	}
}

func (this StorageProxy) Get(key key.Key) (interface{}, error) {
	object, err := this.PreferedStorage.Get(key)
	if err != nil {
		object, err = this.BackupStorage.Get(key)
		if err == nil {
			return object, nil
		} else {
			return nil, err
		}
	}
	if object == nil {
		object, err = this.BackupStorage.Get(key)
		if err != nil {
			return nil, err
		}
		if object != nil {
			this.PreferedStorage.Set(key, object)
		}
	}
	return object, nil
}

func (this StorageProxy) Set(key key.Key, object interface{}) error {
	if object != nil {
		err := this.PreferedStorage.Set(key, object)
		if err != nil {
			return err
		}
		err = this.BackupStorage.Set(key, object)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this StorageProxy) Add(key key.Key, object interface{}) error {
	if object != nil {
		err := this.PreferedStorage.Add(key, object)
		if err != nil {
			return err
		}
		err = this.BackupStorage.Add(key, object)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this StorageProxy) MultiGet(keys []key.Key) (map[key.Key]interface{}, error) {
	resultMap, err := this.PreferedStorage.MultiGet(keys)
	if err != nil {
		return nil, err
	}
	missedKeyCount := 0
	for _, key := range keys {
		if _, find := resultMap[key]; !find {
			missedKeyCount++
		}
	}
	if missedKeyCount > 0 {
		missedKeys := make([]key.Key, missedKeyCount)
		i := 0
		for _, key := range keys {
			if _, find := resultMap[key]; !find {
				missedKeys[i] = key
				i++
			}
		}
		missedMap, err := this.BackupStorage.MultiGet(missedKeys)
		if err != nil {
			return nil, err
		}
		this.PreferedStorage.MultiSet(missedMap)
		for k, v := range missedMap {
			resultMap[k] = v
		}
	}
	return resultMap, nil
}

func (this StorageProxy) MultiSet(objectMap map[key.Key]interface{}) error {
	err := this.PreferedStorage.MultiSet(objectMap)
	if err != nil {
		return err
	}
	err = this.BackupStorage.MultiSet(objectMap)
	if err != nil {
		return err
	}
	return nil
}

func (this StorageProxy) Delete(key key.Key) error {
	err := this.BackupStorage.Delete(key)
	if err != nil {
		return err
	}
	err = this.PreferedStorage.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (this StorageProxy) Incr(key key.Key, step uint64) (newValue uint64, err error) {
	result, err := this.PreferedStorage.(CounterStorage).Incr(key, step)
	if err != nil {
		return result, err
	}
	result, err = this.BackupStorage.(CounterStorage).Incr(key, step)
	if err != nil {
		return result, err
	}
	return result, err
}

func (this StorageProxy) Decr(key key.Key, step uint64) (newValue uint64, err error) {
	result, err := this.PreferedStorage.(CounterStorage).Decr(key, step)
	if err != nil {
		return result, err
	}
	result, err = this.BackupStorage.(CounterStorage).Decr(key, step)
	if err != nil {
		return result, err
	}
	return result, err
}

func (this StorageProxy) FlushAll() {
	this.PreferedStorage.FlushAll()
	this.BackupStorage.FlushAll()
}
