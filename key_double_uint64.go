package storage

import (
	"fmt"
)

type KeyDoubleUint64 struct {
	key1 uint64
	key2 uint64
}

func NewKeyDoubleUint64(key1, key2 uint64) KeyDoubleUint64 {
	return KeyDoubleUint64{key1: key1, key2: key2}
}
func (u KeyDoubleUint64) Equal(u2 KeyDoubleUint64) bool {
	return u.key1 == u2.key1 && u.key2 == u2.key2
}

func (u *KeyDoubleUint64) FromString(str string) bool {
	_, err := fmt.Sscanf(str, "%d-%d", &u.key1, &u.key2)
	return err == nil
}

func (u KeyDoubleUint64) ToString() string {
	return fmt.Sprintf("%d-%d", u.key1, u.key2)
}
func (u KeyDoubleUint64) GetKeys() (key1, key2 uint64) {
	key1 = u.key1
	key2 = u.key2
	return
}
func (u KeyDoubleUint64) GetKey1() (key1 uint64) {
	key1 = u.key1
	return
}
func (u KeyDoubleUint64) GetKey2() (key2 uint64) {
	key2 = u.key2
	return
}
func (u KeyDoubleUint64) GetOtherKey(myUin uint64) (ret uint64, ok bool) {
	if myUin == u.key1 {
		ok = true
		ret = u.key2
		return
	} else if myUin == u.key2 {
		ok = true
		ret = u.key1
		return
	}

	return
}
func (u *KeyDoubleUint64) SetKeys(key1, key2 uint64) {
	u.key1 = key1
	u.key2 = key2
}

type KeyDoubleUint64List []KeyDoubleUint64

func (keys KeyDoubleUint64List) GetOtherKeyList(myUin uint64) (retList []uint64) {
	var otherKey uint64
	var ok bool
	for _, u := range keys {
		otherKey, ok = u.GetOtherKey(myUin)
		if ok {
			retList = append(retList, otherKey)
		}
	}
	return
}
func (keys KeyDoubleUint64List) Remove(Key KeyDoubleUint64) (ret KeyDoubleUint64List) {
	ret = make(KeyDoubleUint64List, len(keys), len(keys))
	var index int = 0
	for _, value := range keys {
		if !value.Equal(Key) {
			ret[index] = value
			index++
		}
	}
	ret = ret[:index]
	return

}

func (keys KeyDoubleUint64List) IsInList(Key KeyDoubleUint64) bool {
	for _, tmpKey := range keys {
		if tmpKey.Equal(Key) {
			return true
		}
	}
	return false
}
