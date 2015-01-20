package storage

import (
	"encoding/json"
	"github.com/0studio/storage_key"
	"github.com/dropbox/godropbox/memcache"
	"github.com/golang/groupcache"
	"reflect"
	"testing"
)

func TestGetSetGC(t *testing.T) {
	tt := T{1}

	jsonEncoding := JsonEncoding{reflect.TypeOf(&tt)}
	client := memcache.NewMockClient()
	mcStorage := NewMcStorage(client, "test", 0, jsonEncoding)
	mcStorage.Set(key.String("1"), tt)
	res, _ := mcStorage.Get(key.String("1"))
	defer mcStorage.Delete(key.String("1"))
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes := res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}

	var groupcache = groupcache.NewGroup("SlowDBCache", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, Key string, dest groupcache.Sink) error {
			result, err := mcStorage.Get(key.String(Key))
			if err != nil {
				return nil
			}
			bytes, err := json.Marshal(result)
			if err != nil {
				return nil
			}
			dest.SetBytes(bytes)
			return nil
		}))
	gcStorage := &GroupCacheKvStorage{groupcache, 0, jsonEncoding}
	res, _ = gcStorage.Get(key.String("1"))
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes = res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}

	mcStorage.Delete(key.String("1"))
	res, _ = gcStorage.Get(key.String("1"))
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes = res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}
}
