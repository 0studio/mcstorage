package storage

import (
	"testing"
	"reflect"
	"sort"
)

func TestGetLimit(t *testing.T) {
	var array []int
	for i := 1; i <= 200; i++ {
		array=append(array,i)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(array)))
	slice:=IntReversedSlice(array)

	jsonEncoding:=JsonEncoding{reflect.TypeOf(&slice)}
	mcStorage := NewMcStorage([]string{"localhost:12000"}, "test", 0, jsonEncoding)

	mcStorage.Set("1", slice)
	result, _ := mcStorage.Getlimit("1",0,0,1,20)
	defer mcStorage.Delete("1")
	if result.(IntReversedSlice).Len()!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=200{
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19]!=181{
		t.Error("first one should be 181")
	}

	result, _ = mcStorage.Getlimit("1",0,200,1,20)
	if result.(IntReversedSlice).Len()!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=199{
		t.Error("first one should be 199")
	}
	if result.(IntReversedSlice)[19]!=180{
		t.Error("first one should be 180")
	}

}

func TestAddItem(t *testing.T) {
	var array []int
	for i := 1; i <= 200; i++ {
		array=append(array,i)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(array)))
	slice:=IntReversedSlice(array)
	jsonEncoding:=JsonEncoding{reflect.TypeOf(&slice)}
	mcStorage := NewMcStorage([]string{"localhost:12000"}, "test", 0, jsonEncoding)
	mcStorage.Set("1", slice)
	result, _ := mcStorage.Getlimit("1",0,0,1,20)
	defer mcStorage.Delete("1")
	if result.(IntReversedSlice).Len()!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=200{
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19]!=181{
		t.Error("first one should be 181")
	}

	mcStorage.AddItem("1",201)
	result, _ = mcStorage.Getlimit("1",0,0,1,20)
	if result.(IntReversedSlice).Len()!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=201{
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19]!=182{
		t.Error("first one should be 182")
	}

	mcStorage.DeleteItem("1",193)
	result, _ = mcStorage.Getlimit("1",0,0,1,20)
	if result.(IntReversedSlice).Len()!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=201{
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19]!=181{
		t.Error("first one should be 181")
	}

}