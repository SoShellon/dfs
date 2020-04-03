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
	dirRef       map[string]int
}

var s *SingleFileSystem

//InitStorageNode initialize the core system of storage server
func InitStorageNode(clientPort int, commandPort int, namingServer int, rootPath string) error {
	err := checkAndCreateRoot(rootPath)
	if err != nil {
		return err
	}

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
		computeDirectoryReference(files),
	}
	registerData := s.StorageNode
	resp := struct {
		Files         []string `json:"files"`
		ExceptionType string   `json:"exception_type"`
		ExceptionInfo string   `json:"exception_info"`
	}{}
	err = common.SendRequest(fmt.Sprintf("localhost:%d/register", namingServer), registerData, &resp)
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

func getAncestors(file string) []string {
	dirs := []string{}
	idx := strings.LastIndex(file, "/")
	for idx > 0 {
		dir := file[0:idx]
		dirs = append(dirs, dir)
		file = dir
		idx = strings.LastIndex(file, "/")
	}
	return dirs
}
func computeDirectoryReference(files []string) map[string]int {
	ref := map[string]int{}
	for _, file := range files {
		dirs := getAncestors(file)
		for _, dir := range dirs {
			if count, exists := ref[dir]; !exists {
				ref[dir] = 1
			} else {
				ref[dir] = count + 1
			}

		}
	}
	return ref
}

func checkAndCreateRoot(root string) error {
	_, err := os.Stat(root)
	if os.IsNotExist(err) {
		err = os.MkdirAll(root, 0777)
		return err
	}
	return err
}
func walkThroughRootPath(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		path = strings.TrimPrefix(path, root)
		if path == "" || info.IsDir() {
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
	if err := os.MkdirAll(filepath.Dir(full), 0770); err != nil {
		return false, err
	}
	_, err := os.Create(full)
	if err != nil {
		return false, err
	}
	s.Files = append(s.Files, path)
	return true, nil
}

//DeleteFile deletes a file on the file system
func (s *SingleFileSystem) DeleteFile(path string) (bool, error) {
	if path == "/" {
		return false, nil
	}
	full := s.getFullPath(path)
	fileInfo, err := os.Stat(full)
	if err != nil {
		return false, err
	}
	deleteFiles := []string{}
	if fileInfo.IsDir() {
		all, err := walkThroughRootPath(full)
		if err != nil {
			return false, err
		}
		for _, child := range all {
			deleteFiles = append(deleteFiles, path+child)
		}
	} else {
		idx := s.indexFile(path)
		if idx >= 0 {
			deleteFiles = append(deleteFiles, path)
		}
	}
	if len(deleteFiles) == 0 {
		return false, errors.New("no such file")
	}
	for _, path = range deleteFiles {
		full = s.getFullPath(path)
		err = os.Remove(full)
		if err != nil {
			return false, err
		}

		//prune
		dirs := getAncestors(path)
		for _, dir := range dirs {
			count := s.dirRef[dir]
			count--
			if count == 0 {
				delete(s.dirRef, dir)
				os.Remove(s.getFullPath(dir))
			} else {
				s.dirRef[dir] = count
			}
		}

		idx := s.indexFile(path)
		if idx == -1 {
			return false, errors.New("no such file")
		}
		s.Files[len(s.Files)-1], s.Files[idx] = s.Files[idx], s.Files[len(s.Files)-1]
		s.Files = s.Files[:len(s.Files)-1]
	}
	return true, nil
}

//CopyFile copy a file from another storage server
func (s *SingleFileSystem) CopyFile(path string, node *common.StorageNode) (bool, error) {
	size, err := s.getRemoteSize(path, node)
	if err != nil {
		return false, err
	}
	_, err = s.CreateFile(path)
	if err != nil {
		return false, err
	}
	data, err := s.getRemoteBytes(path, size, node)
	if err != nil {
		return false, err
	}
	_, err = s.Write(path, 0, data)
	if err != nil {
		return false, err
	}
	fullName := s.getFullPath(path)
	err = os.Truncate(fullName, int64(size))
	return true, err
}
func (s *SingleFileSystem) getRemoteBytes(path string, size int, storageNode *common.StorageNode) ([]byte, error) {
	req := &struct {
		Path   string `json:"path"`
		Offset int    `json:"offset"`
		Length int    `json:"length"`
	}{path, 0, size}
	resp := &struct {
		Data          []byte `json:"data"`
		ExceptionInfo string `json:"exception_info"`
	}{}
	common.SendRequest(fmt.Sprintf("%s:%d/storage_read", storageNode.StorageIP, storageNode.ClientPort), req, resp)
	if len(resp.ExceptionInfo) > 0 {
		return nil, errors.New(resp.ExceptionInfo)
	}
	return resp.Data, nil
}
func (s *SingleFileSystem) getRemoteSize(path string, storageNode *common.StorageNode) (int, error) {
	req := &struct {
		Path string `json:"path"`
	}{path}
	resp := &struct {
		Size          int    `json:"size"`
		ExceptionInfo string `json:"exception_info"`
	}{}
	common.SendRequest(fmt.Sprintf("%s:%d/storage_size", storageNode.StorageIP, storageNode.ClientPort), req, resp)
	if len(resp.ExceptionInfo) > 0 {
		return 0, errors.New(resp.ExceptionInfo)
	}
	return resp.Size, nil
}

//ValidatePath determines whether the path is valid
func (s *SingleFileSystem) ValidatePath(path string) bool {
	return len(path) > 0
}

//GetFileSize returns the size of a file
func (s *SingleFileSystem) GetFileSize(path string) (int64, error) {
	full := s.getFullPath(path)
	fileInfo, err := os.Stat(full)
	if err != nil {
		return 0, err
	}
	if fileInfo.IsDir() {
		return 0, errors.New("is not file")
	}
	return fileInfo.Size(), nil
}

//WithinBounds determine params is within the bounds of file
func (s *SingleFileSystem) WithinBounds(path string, offset int64, length int64) bool {
	if offset < 0 || length < 0 {
		return false
	}
	size, err := s.GetFileSize(path)
	if err != nil {
		//fall through
		return true
	}
	return size >= offset+length
}

//Read return the bytes of specified file
func (s *SingleFileSystem) Read(path string, offset int64, length int64) ([]byte, error) {
	_, err := s.GetFileSize(path)
	if err != nil {
		return nil, err
	}
	res := make([]byte, length)
	f, err := os.Open(s.getFullPath(path))
	defer f.Close()
	if err != nil {
		return nil, err
	}
	_, err = f.Seek(offset, 0)
	if err != nil {
		return nil, err
	}
	_, err = f.Read(res)
	return res, err
}

//Write writes bytes into a file
func (s *SingleFileSystem) Write(path string, offset int64, data []byte) (bool, error) {
	_, err := s.GetFileSize(path)
	if err != nil {
		return false, err
	}
	f, err := os.OpenFile(s.getFullPath(path), os.O_WRONLY, os.ModeAppend)
	defer f.Close()
	if err != nil {
		return false, err
	}
	_, err = f.WriteAt(data, offset)
	if err != nil {
		return false, err
	}
	err = f.Sync()
	return err == nil, err
}
