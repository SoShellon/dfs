package core

import (
	"testing"

	"cmu.edu/dfs/common"
)

func TestRegister(t *testing.T) {
	InitRegistrar()
	r := GetRegistrar()
	file1 := "/a"
	root := "/"
	r.AddStorageNode(&common.StorageNode{Files: []string{file1}})
	list, _ := r.ListFiles(root)

	if len(list) != 1 || list[0] != "/a" {
		t.Errorf("should be 1 file:%+v", list)
	}
	file2 := "/b"
	_, err := r.AddStorageNode(&common.StorageNode{Files: []string{file2}})
	if err == nil {
		t.Error("should not add this second node")
	}
	list, err = r.AddStorageNode(&common.StorageNode{CommandPort: 1223, Files: []string{file1, file2}})
	if len(list) != 1 {
		t.Errorf("should be 1 duplicated file:%+v, err:%+v", list, err)
	}
	list, _ = r.ListFiles(root)
	if len(list) != 2 {
		t.Errorf("should be 2 files:%+v", list)
	}
}

func TestValidate(t *testing.T) {
	InitRegistrar()
	r := GetRegistrar()
	if r.ValidatePath("") {
		t.Error("should be invalid")
	}
}

func TestList(t *testing.T) {
	InitRegistrar()
	r := GetRegistrar()
	files, err := r.ListFiles("/")
	if err != nil {
		t.Errorf("should be no error:%+v", err)
	}
	files, err = r.ListFiles("/file1")
	if err == nil {
		t.Errorf("should be error:%+v", files)
	}
}

func TestFileNodeSearch(t *testing.T) {
	InitRegistrar()
	r := GetRegistrar()
	r.AddStorageNode(&common.StorageNode{Files: []string{"/a/b"}})
	succ, err := r.CreateDir("/a/b")
	if err != nil {
		t.Errorf("should be no error:%v", err)
	}
	if succ {
		t.Error("should fail to create")
	}
	parts := []string{"a", "b"}
	fileNode := r.getParentNode(parts)
	if fileNode.search(parts[1:]) == nil {
		t.Errorf("should has child:%+v", fileNode)
	}
}

func TestCreateFile(t *testing.T) {
	InitRegistrar()
	r := GetRegistrar()
	r.CreateDir("/directory/subdirectory")
	r.CreateDir("/another_directory")
	succ, err := r.CreateFile("//file")
	if !succ || err != nil {
		t.Errorf("should be success:%v %v", succ, err)
	}
}
