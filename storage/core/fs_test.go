package core

import (
	"strings"
	"testing"
)

func TestWalkThroughDir(t *testing.T){
	files, err:=walkThroughRootPath("../../API")
	if err!=nil {
		t.Error(err.Error())
	}else if len(files)!=4 {
		t.Errorf("API dir walkthrough length error: %+v", files)
	}
	if !strings.HasPrefix(files[0],"/API") {
		t.Errorf("path invalid:%s",files[0])
	}
}

func TestCreateAndDeleteFile(t *testing.T) {
	s := &SingleFileSystem{rootPath:"../../build"}
	s.Files = []string{}
	path:="test"
	success, err:=s.DeleteFile(path)
	if success|| err==nil {
		t.Error("should be error")
	}
	success, err= s.CreateFile(path)
	if !success || err!=nil {
		t.Errorf("invalid %v %s",success, err.Error())
	}
	if s.getFullPath(s.Files[0])!=s.rootPath+path{
		t.Errorf("wrong files:%+v", s.rootPath+path)
	}
	success, err = s.CreateFile(path)
	if success{
		t.Errorf("should be error")
	}
	success, err= s.DeleteFile(path)
	if !success {
		t.Errorf("should be deleted")
	}
	s = nil
}