package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cmu.edu/dfs/common"
)

//SingleFileSystem manage the file sytem of a storage server
type SingleFileSystem struct {
	common.StorageNode
	rootPath     string
	namingServer int
}

var s *SingleFileSystem

//InitStorageNode initialize the core system of storage server
func InitStorageNode(clientPort int, commandPort int, namingServer int, rootPath string) error {

	files, err := walkThroughRootPath(rootPath)
	if err != nil {
		return err
	}
	s = &SingleFileSystem{
		common.StorageNode{
			CommandPort: commandPort,
			ClientPort:  clientPort,
			StorageIP:   "localhost",
			Files:       files,
		},
		rootPath,
		namingServer,
	}
	registerData := s.StorageNode
	resp := struct {
		Files         []string `json:"files"`
		ExceptionType string   `json:"exception_type"`
		ExceptionInfo string   `json:"exception_info"`
	}{}
	err = common.SendRequest(fmt.Sprintf("localhost:%d/registration", namingServer), registerData, &resp)
	if err != nil {
		return err
	}
	if resp.ExceptionInfo != "" {
		common.Error(resp.ExceptionInfo)
		return errors.New("registration fails")
	}
	for _, file := range resp.Files {
		s.DeleteFile(file)
	}
	return nil
}

func walkThroughRootPath(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		path = strings.TrimPrefix(path, root)
		if path == "" {
			return nil
		}
		if strings.HasPrefix(path, "/") {
			files = append(files, path)
		} else {
			files = append(files, "/"+path)
		}
		return nil
	})
	return
}

//GetStorageNode return the instance of file system
func GetStorageNode() *SingleFileSystem {
	return s
}
func (s *SingleFileSystem) getFullPath(path string) string {
	return s.rootPath + path
}

func (s *SingleFileSystem) indexFile(path string) int {
	for i, file := range s.Files {
		if file == path {
			return i
		}
	}
	return -1
}

//CreateFile creates a file on the file system
func (s *SingleFileSystem) CreateFile(path string) (bool, error) {
	if s.indexFile(path) != -1 {
		return false, nil
	}
	full := s.getFullPath(path)
	_, err := os.Create(full)
	if err != nil {
		return false, err
	}
	s.Files = append(s.Files, path)
	return true, nil
}

//DeleteFile deletes a file on the file system
func (s *SingleFileSystem) DeleteFile(path string) (bool, error) {
	idx := s.indexFile(path)
	if idx == -1 {
		return false, errors.New("no such file")
	}
	full := s.getFullPath(path)
	err := os.Remove(full)
	if err != nil {
		return false, err
	}
	s.Files[len(s.Files)-1], s.Files[idx] = s.Files[idx], s.Files[len(s.Files)-1]
	s.Files = s.Files[:len(s.Files)-1]
	return true, nil
}

//CopyFile copy a file from another storage server
func (s *SingleFileSystem) CopyFile(path string, node *common.StorageNode) (bool, error) {
	return false, nil
}

//GetFileSize returns the size of a file
func (s *SingleFileSystem) GetFileSize(path string) (int64, error) {
	full := s.getFullPath(path)
	fileInfo, err := os.Stat(full)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

//Read return the bytes of specified file
func (s *SingleFileSystem) Read(path string, offset int64, length int64) ([]byte, error) {
	return nil, nil
}

//Write writes bytes into a file
func (s *SingleFileSystem) Write(path string, offset int64, data []byte) (bool, error) {
	return false, nil
}
