package storage

import (
	"github.com/0studio/storage_key"
	"github.com/dropbox/godropbox/memcache"
	"reflect"
	"sort"
	"testing"
)

func TestGetLimit(t *testing.T) {
	var array []int
	for i := 1; i <= 200; i++ {
		array = append(array, i)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(array)))
	slice := IntReversedSlice(array)

	jsonEncoding := JsonEncoding{reflect.TypeOf(&slice)}

	client := memcache.NewMockClient()
	mcStorage := NewMcStorage(client, "test", 0, jsonEncoding)

	mcStorage.Set(key.String("1"), slice)
	result, _ := mcStorage.Getlimit(key.String("1"), 0, 0, 1, 20)
	defer mcStorage.Delete(key.String("1"))
	if result.(IntReversedSlice).Len() != 20 {
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0] != 200 {
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19] != 181 {
		t.Error("first one should be 181")
	}

	result, _ = mcStorage.Getlimit(key.String("1"), 0, 200, 1, 20)
	if result.(IntReversedSlice).Len() != 20 {
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0] != 199 {
		t.Error("first one should be 199")
	}
	if result.(IntReversedSlice)[19] != 180 {
		t.Error("first one should be 180")
	}

}

func TestAddItem(t *testing.T) {
	var array []int
	for i := 1; i <= 200; i++ {
		array = append(array, i)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(array)))
	slice := IntReversedSlice(array)
	jsonEncoding := JsonEncoding{reflect.TypeOf(&slice)}

	client := memcache.NewMockClient()
	mcStorage := NewMcStorage(client, "test", 0, jsonEncoding)

	mcStorage.Set(key.String("1"), slice)
	result, _ := mcStorage.Getlimit(key.String("1"), 0, 0, 1, 20)
	defer mcStorage.Delete(key.String("1"))
	if result.(IntReversedSlice).Len() != 20 {
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0] != 200 {
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19] != 181 {
		t.Error("first one should be 181")
	}

	mcStorage.AddItem(key.String("1"), 201)
	result, _ = mcStorage.Getlimit(key.String("1"), 0, 0, 1, 20)
	if result.(IntReversedSlice).Len() != 20 {
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0] != 201 {
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19] != 182 {
		t.Error("first one should be 182")
	}

	mcStorage.DeleteItem(key.String("1"), 193)
	result, _ = mcStorage.Getlimit(key.String("1"), 0, 0, 1, 20)
	if result.(IntReversedSlice).Len() != 20 {
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0] != 201 {
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19] != 181 {
		t.Error("first one should be 181")
	}

}
