package core

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"cmu.edu/dfs/common"
)

const (
	errorNotDir   = "no such dir"
	errorNoFile   = "no such file"
	errorNoParent = "no parent directory"
)

type fileNode struct {
	isDir      bool
	childNodes map[string]*fileNode
	index      string
	slave      []string
	token      string
	rwlock     sync.RWMutex
	countLock  sync.Mutex
	readCount  int32
	parent     *fileNode
}

func (f *fileNode) lock(exclusive bool) {
	if exclusive {
		f.rwlock.Lock()
	} else {
		f.rwlock.RLock()
	}
}

func (f *fileNode) unlock(exclusive bool) {
	if exclusive {
		f.rwlock.Unlock()
	} else {
		f.rwlock.RUnlock()
	}
}

func (f *fileNode) getFullName() string {
	res := ""
	cur := f
	for cur.parent != nil {
		res = "/" + cur.token + res
		cur = cur.parent
	}
	return res
}

func (f *fileNode) delete() bool {
	if f.parent == nil {
		return false
	}
	delete(f.parent.childNodes, f.token)
	return true
}

func (f *fileNode) search(parts []string) *fileNode {
	cur := f
	for _, part := range parts {
		child := cur.childNodes[part]
		if child == nil {
			return nil
		}
		cur = child
	}
	return cur
}

func (f *fileNode) validate(parts []string) bool {
	cur := f
	for _, part := range parts {
		if !cur.isDir {
			return false
		}
		next := cur.childNodes[part]
		if next == nil {
			return true
		}
		cur = next
	}
	return false
}

func (f *fileNode) insert(parts []string, index string, isDir bool) {
	cur := f
	for _, part := range parts {
		if cur.childNodes[part] == nil {
			node := buildDir()
			node.token = part
			node.parent = cur
			node.index = index
			cur.childNodes[part] = node
		}
		cur = cur.childNodes[part]
	}
	cur.isDir = isDir
	cur.index = index
	common.Log("%s add new file node:%+v", index, parts)
}
func buildDir() *fileNode {
	return &fileNode{
		isDir:      true,
		childNodes: map[string]*fileNode{},
		rwlock:     sync.RWMutex{},
		readCount:  0,
	}
}

//Registrar manages registered storage nodes
type Registrar struct {
	storageNodes map[string]*common.StorageNode
	globalRoot   *fileNode
}

var r *Registrar

//InitRegistrar initialize the core system
func InitRegistrar() {
	root := buildDir()
	r = &Registrar{
		globalRoot:   root,
		storageNodes: make(map[string]*common.StorageNode),
	}
}

//GetRegistrar expose the unique registrar instance
func GetRegistrar() *Registrar {
	return r
}

//ValidatePath determine a path is valid
func (r *Registrar) ValidatePath(path string) bool {
	if len(path) == 0 {
		return false
	}
	return true
}

//AddStorageNode tries to add a new storage node
func (r *Registrar) AddStorageNode(node *common.StorageNode) ([]string, error) {
	key := node.GetIndexKey()
	if r.storageNodes[key] != nil {
		return nil, errors.New("this storage client already registered:" + key)
	}
	r.storageNodes[key] = node
	duplicates := make([]string, 0)
	// for _, path := range node.Files {
	// 	succ, _ := r.CreateFile(path)
	// 	if !succ {
	// 		duplicates = append(duplicates, path)
	// 	}
	// }
	for i, file := range node.GetFileTokens() {
		valid := r.globalRoot.validate(file)
		if valid {
			r.globalRoot.insert(file, key, false)
		} else if node.Files[i] != "/" {
			duplicates = append(duplicates, node.Files[i])
		}
	}
	return duplicates, nil
}

func (r *Registrar) getStorageNodeWithIndex(index string) *common.StorageNode {
	return r.storageNodes[index]
}
func (r *Registrar) getStorageNode(node *fileNode) *common.StorageNode {
	//root node
	if node.parent == nil {
		for _, node := range r.storageNodes {
			return node
		}
	}
	return r.getStorageNodeWithIndex(node.index)
}

//Exists checks whether a file exists
func (r *Registrar) Exists(file string) bool {
	if !strings.HasPrefix(file, "/") {
		return false
	}
	tokens := common.Tokenize(file)
	return r.globalRoot.search(tokens) != nil
}

//GetStorageNode find the storage node according to the path
func (r *Registrar) GetStorageNode(path string) (*common.StorageNode, error) {
	parts := common.Tokenize(path)
	node := r.globalRoot.search(parts)
	if node == nil {
		return nil, errors.New(errorNoFile)
	}
	return r.storageNodes[node.index], nil
}

func (r *Registrar) getParentNode(parts []string) *fileNode {
	if len(parts) == 0 {
		return nil
	}
	parentParts := parts[:len(parts)-1]
	fileNode := r.globalRoot.search(parentParts)
	return fileNode
}

//CreateFile create a file in the file system
func (r *Registrar) CreateFile(path string) (bool, error) {
	parts := common.Tokenize(path)
	if len(parts) == 0 {
		return false, nil
	}
	fileNode := r.getParentNode(parts)
	if fileNode == nil {
		return false, errors.New(errorNoParent)
	}
	if !fileNode.isDir {
		return false, errors.New(errorNotDir)
	}
	sub := parts[len(parts)-1:]
	child := fileNode.search(sub)
	if child != nil {
		return false, nil
	}
	storageNode := r.getStorageNode(fileNode)
	req := &struct {
		Path string `json:"path"`
	}{path}
	resp := &struct {
		Success       bool   `json:"success"`
		ExceptionInfo string `json:"exception_info"`
	}{}
	err := common.SendRequest(fmt.Sprintf("%s:%d/storage_create", storageNode.StorageIP, storageNode.CommandPort), req, resp)
	if err != nil {
		return false, err
	}
	if !resp.Success {
		if len(resp.ExceptionInfo) > 0 {
			return false, errors.New(resp.ExceptionInfo)
		}
		return false, nil
	}
	fileNode.insert(sub, storageNode.GetIndexKey(), false)
	return true, nil
}

//CreateDir creates a directory in the file system
func (r *Registrar) CreateDir(path string) (bool, error) {
	parts := common.Tokenize(path)
	if len(parts) == 0 {
		return false, nil
	}
	fileNode := r.getParentNode(parts)
	if fileNode == nil {
		return false, errors.New(errorNoParent)
	}
	if !fileNode.isDir {
		return false, errors.New(errorNotDir)
	}
	sub := parts[len(parts)-1:]
	child := fileNode.search(sub)
	if child != nil {
		return false, nil
	}
	storageNode := r.getStorageNode(fileNode)
	fileNode.insert(sub, storageNode.GetIndexKey(), true)
	return true, nil
}

//Delete delete a file or dir from the file system
func (r *Registrar) Delete(path string) (bool, error) {
	parts := common.Tokenize(path)
	fileNode := r.globalRoot.search(parts)
	if fileNode == nil {
		return false, errors.New(errorNoFile)
	}
	list := []string{}
	for key := range r.storageNodes {
		list = append(list, key)
	}
	r.deleteReplications(path, list)
	return fileNode.delete(), nil
}

//ListFiles list all the files belonging to a dir
func (r *Registrar) ListFiles(path string) ([]string, error) {
	parts := common.Tokenize(path)
	fileNode := r.globalRoot.search(parts)
	if fileNode == nil {
		return nil, errors.New(errorNoFile)
	}
	res := []string{}
	prefix := fileNode.getFullName()
	if !fileNode.isDir {
		res = append(res, prefix)
		return res, nil
	}
	for _, node := range fileNode.childNodes {
		res = append(res, node.token)
	}
	return res, nil
}

//IsDir determine whether a file is a directory
func (r *Registrar) IsDir(path string) (bool, error) {
	parts := common.Tokenize(path)
	fileNode := r.globalRoot.search(parts)
	if fileNode == nil {
		return false, errors.New("no file found")
	}
	return fileNode.isDir, nil
}

//Lock locks the file node
func (r *Registrar) Lock(path string, exclusive bool) {
	parts := common.Tokenize(path)
	cur := r.globalRoot
	for _, part := range parts {
		cur.lock(false)
		cur = cur.childNodes[part]
	}
	cur.lock(exclusive)

}

func (r *Registrar) getReplicationNodes(file *fileNode) []string {
	others := []string{}
	for key, value := range r.storageNodes {
		if key != file.index {
			others = append(others, value.GetIndexKey())
		}
	}
	return others
}
func (r *Registrar) deleteReplications(path string, list []string) {
	req := &struct {
		Path string `json:"path"`
	}{path}
	resp := &struct {
		Success       bool   `json:"success"`
		ExceptionInfo string `json:"exception_info"`
	}{}
	for _, index := range list {
		storageNode := r.getStorageNodeWithIndex(index)
		common.SendRequest(fmt.Sprintf("%s:%d/storage_delete", storageNode.StorageIP, storageNode.CommandPort), req, resp)
	}
}

//Unlock unlocks the file node
func (r *Registrar) Unlock(path string, exclusive bool) {
	parts := common.Tokenize(path)
	cur := r.globalRoot
	for _, part := range parts {
		cur.unlock(false)
		cur = cur.childNodes[part]
	}
	cur.unlock(exclusive)
	if !cur.isDir {
		if !exclusive {
			//read request
			if cur.readCount >= 20 {
				cur.countLock.Lock()
				defer cur.countLock.Unlock()
				if cur.readCount >= 20 {
					cur.readCount = 0
				}
				master := r.getStorageNode(cur)
				req := &struct {
					Path       string `json:"path"`
					ServerIP   string `json:"server_ip"`
					ServerPort int    `json:"server_port"`
				}{path, master.StorageIP, master.ClientPort}
				resp := &struct {
					Success       bool   `json:"success"`
					ExceptionInfo string `json:"exception_info"`
				}{}
				cur.slave = r.getReplicationNodes(cur)
				for _, index := range cur.slave {
					storageNode := r.getStorageNodeWithIndex(index)
					common.SendRequest(fmt.Sprintf("%s:%d/storage_copy", storageNode.StorageIP, storageNode.CommandPort), req, resp)
				}
			} else {
				atomic.AddInt32(&cur.readCount, 1)
			}
		} else {
			//write request
			r.deleteReplications(path, cur.slave)
		}
	}
}
