package storage

/**
import (
	"reflect"
	"testing"
	"github.com/dropbox/godropbox/memcache"
)

func TestIncrDecr(t *testing.T) {
	jsonEncoding:=JsonEncoding{reflect.TypeOf(1)}

	client:=memcache.NewMockClient()
	mcStorage := NewMcStorage(client, "test", 0, jsonEncoding)

	mcStorage.Set(String("1"), 1)
	res, _ := mcStorage.Get(String("1"))
	defer mcStorage.Delete(String("1"))
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if res.(int)!=1{
		t.Error("value should be 1")
	}

	resIncr,_:=mcStorage.Incr(String("1"),1,1)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resIncr!=2{
		t.Error("value should be 2")
	}

	resIncr,_=mcStorage.Incr(String("1"),3,2)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resIncr!=5{
		t.Error("value should be 5")
	}

	resDecr,_:=mcStorage.Decr(String("1"),1)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resDecr!=4{
		t.Error("value should be 4")
	}

	resDecr,_=mcStorage.Decr(String("1"),2)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resDecr!=2{
		t.Error("value should be 2")
	}

	resDecr,err:=mcStorage.Decr(String("2"),2)
	if err!=nil{
		t.Error("err should be nil",err)
	}
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resDecr!=0{
		t.Error("value should be 0")
	}

}
*/
