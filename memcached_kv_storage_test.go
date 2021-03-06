package storage

import (
	"github.com/0studio/storage_key"
	"github.com/dropbox/godropbox/memcache"
	"reflect"
	"strconv"
	"testing"
)

type T struct {
	A int
}

type TSlice []T

func TestGetSet(t *testing.T) {
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
}

func TestGetSetNil(t *testing.T) {
	tt := T{1}
	var tt1 T

	jsonEncoding := JsonEncoding{reflect.TypeOf(&tt)}
	client := memcache.NewMockClient()
	mcStorage := NewMcStorage(client, "test", 0, jsonEncoding)
	//mcStorage := NewMcStorage([]string{"localhost:12000"}, "test", 0, jsonEncoding)
	mcStorage.Set(key.String("1"), tt1)
	res, _ := mcStorage.Get(key.String("1"))
	defer mcStorage.Delete(key.String("1"))
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
}

func TestGetSetNilSlice(t *testing.T) {

	var nilSlice TSlice
	jsonEncoding := JsonEncoding{reflect.TypeOf(&nilSlice)}
	//mcStorage := NewMcStorage([]string{"localhost:12000"}, "test", 0, jsonEncoding)

	client := memcache.NewMockClient()
	mcStorage := NewMcStorage(client, "test", 0, jsonEncoding)

	mcStorage.Set(key.String("1"), nilSlice)
	res, err := mcStorage.Get(key.String("1"))
	defer mcStorage.Delete(key.String("1"))
	if res != nil {
		t.Error("result should be nil")
	}
	if err != nil {
		t.Error("result should be nil")
	}
}

func TestMultiGetSet(t *testing.T) {
	tt := T{1}
	jsonEncoding := JsonEncoding{reflect.TypeOf(&tt)}
	///mcStorage := NewMcStorage([]string{"localhost:12000"}, "test", 0, jsonEncoding)

	client := memcache.NewMockClient()
	mcStorage := NewMcStorage(client, "test", 0, jsonEncoding)

	valueMap := make(map[key.Key]interface{})
	keys := make([]key.Key, 10)
	for i := 0; i < 10; i++ {
		keys[i] = key.String(strconv.Itoa(i))
		valueMap[key.String(strconv.Itoa(i))] = T{i}
		defer mcStorage.Delete(key.String(strconv.Itoa(i)))
	}
	mcStorage.MultiSet(valueMap)
	res, _ := mcStorage.MultiGet(keys)
	for k, v := range res {
		if reflect.TypeOf(v) != reflect.TypeOf(tt) {
			t.Error("res type is not T")
		}
		kint, err := strconv.Atoi(k.ToString())
		if err != nil {
			t.Error("key %s is not int ", k)
		}
		vT := v.(T)
		if kint != vT.A {
			t.Error("value should be %s,while it is %s", kint, vT.A)
		}
	}
}

func TestGetSetDelete(t *testing.T) {
	tt := T{1}
	jsonEncoding := JsonEncoding{reflect.TypeOf(&tt)}
	//mcStorage := NewMcStorage([]string{"localhost:12000"}, "test", 0, jsonEncoding)
	client := memcache.NewMockClient()
	mcStorage := NewMcStorage(client, "test", 0, jsonEncoding)

	mcStorage.Set(key.String("1"), tt)
	res, _ := mcStorage.Get(key.String("1"))
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes := res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}
	mcStorage.Delete(key.String("1"))
	res, _ = mcStorage.Get(key.String("1"))
	if res != nil {
		t.Error("res should be nil ,after delete")
	}
}
