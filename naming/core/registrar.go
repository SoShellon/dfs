package core

import (
	"errors"

	"cmu.edu/dfs/common"
)

const (
	errorNoFile   = "no such file"
	errorNoParent = "no parent directory"
)

type fileNode struct {
	isDir      bool
	childNodes map[string]*fileNode
	index      string
	token      string
	parent     *fileNode
}

func (f *fileNode) getFullName() string {
	res := ""
	cur := f
	for cur.parent != nil {
		res = cur.token + res
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

func (f *fileNode) insert(parts []string, index string) {
	cur := f
	for _, part := range parts {
		if cur.childNodes[part] == nil {
			node := buildDir()
			node.token = part
			node.parent = cur
			cur.childNodes[part] = node
		}
		cur = cur.childNodes[part]
	}
	cur.isDir = false
	cur.index = index
}
func buildDir() *fileNode {
	return &fileNode{
		isDir:      true,
		childNodes: map[string]*fileNode{},
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

//AddStorageNode tries to add a new storage node
func (r *Registrar) AddStorageNode(node *common.StorageNode) ([]string, error) {
	key := node.GetIndexKey()
	if r.storageNodes[key] != nil {
		return nil, errors.New("this storage client already registered")
	}
	r.storageNodes[key] = node
	duplicates := []string{}
	for i, file := range node.GetFileTokens() {
		valid := r.globalRoot.validate(file)
		if valid {
			r.globalRoot.insert(file, key)
		} else {
			duplicates = append(duplicates, node.Files[i])
		}
	}
	return duplicates, nil
}

//Exists checks whether a file exists
func (r *Registrar) Exists(file string) bool {
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
	parentParts := parts[0:len(parts)]
	fileNode := r.globalRoot.search(parentParts)
	return fileNode
}

//CreateFile create a file in the file system
func (r *Registrar) CreateFile(path string) (bool, error) {
	parts := common.Tokenize(path)
	parentDir := r.getParentNode(parts)
	if parentDir != nil {
		return false, errors.New(errorNoParent)
	}
	return true, nil
}

//CreateDir creates a directory in the file system
func (r *Registrar) CreateDir(path string) (bool, error) {
	parts := common.Tokenize(path)
	fileNode := r.getParentNode(parts)
	if fileNode != nil {
		return false, errors.New(errorNoParent)
	}
	return true, nil
}

//Delete delete a file or dir from the file system
func (r *Registrar) Delete(path string) (bool, error) {
	parts := common.Tokenize(path)
	fileNode := r.globalRoot.search(parts)
	if fileNode == nil {
		return false, errors.New(errorNoFile)
	}
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
		res = append(res, prefix+"/"+node.token)
	}
	return res, nil
}
