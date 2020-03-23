package core

import (
	"testing"

	"cmu.edu/dfs/common"
)


func TestRegister(t *testing.T) {
	InitRegistrar()
	r:= GetRegistrar()
	file1:="/a"
	root :="/"
	r.AddStorageNode(&common.StorageNode{Files:[]string{file1}})
	list, _:=r.ListFiles(root)

	if len(list)!=1||list[0]!="/a" {
		t.Errorf("should be 1 file:%+v", list)
	}
	file2:="/b"
	_, err:= r.AddStorageNode(&common.StorageNode{Files:[]string{file2}})
	if err==nil {
		t.Error("should not add this second node")
	}
	list,err=r.AddStorageNode(&common.StorageNode{CommandPort:1223, Files:[]string{file1,file2}})
	if len(list)!=1 {
		t.Errorf("should be 1 duplicated file:%+v, err:%+v", list, err)
	}
	list,_ = r.ListFiles(root)
	if len(list)!=2 {
		t.Errorf("should be 2 files:%+v", list)
	}
}
