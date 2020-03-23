package common

import "testing"

func TestGetFileTokens(t *testing.T) {
	node := &StorageNode{
		Files: []string{"/file", "/dir/file"},
	}
	fileTokens := node.GetFileTokens()
	if len(fileTokens) != 2 {
		t.Errorf("file size is %d, not 2", len(fileTokens))
	}

	if len(fileTokens[0]) != 1 || fileTokens[0][0] != "file" {
		t.Errorf("token %+v", fileTokens[0])
	}

	if len(fileTokens[1]) != 2 || fileTokens[1][1] != "file" {
		t.Errorf("token %+v", fileTokens[1])
	}
}
