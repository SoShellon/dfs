package core

import "errors"

type fileNode struct {
	isDir      bool
	childNodes map[string]*fileNode
}

func (f *fileNode) validate(parts []string) bool {
	cur := f
	for _, part := range parts {
		if !cur.isDir {
			return false
		}
		cur = cur.childNodes[part]
		if cur == nil {
			return cur.isDir
		}
	}
	return !cur.isDir
}

func (f *fileNode) insert(parts []string) {
	cur := f
	for _, part := range parts {
		if cur.childNodes[part] == nil {
			cur.childNodes[part] = buildDir()
		}
		cur = cur.childNodes[part]
	}
	cur.isDir = false
}
func buildDir() *fileNode {
	return &fileNode{true, map[string]*fileNode{}}
}

//Registrar manages registered storage nodes
type Registrar struct {
	storageNodes map[string]*StorageNode
	globalRoot   *fileNode
}

var r *Registrar

//InitRegistrar initialize the core system
func InitRegistrar() {
	root := buildDir()
	r = &Registrar{
		globalRoot:   root,
		storageNodes: make(map[string]*StorageNode),
	}
}

//GetRegistrar expose the unique registrar instance
func GetRegistrar() *Registrar {
	return r
}

//AddStorageNode tries to add a new storage node
func (r *Registrar) AddStorageNode(node *StorageNode) ([]string, error) {
	key := node.GetIndexKey()
	if r.storageNodes[key] != nil {
		return nil, errors.New("this storage client already registered")
	}
	r.storageNodes[key] = node
	duplicates := []string{}
	for i, file := range node.GetFileTokens() {
		valid := r.globalRoot.validate(file)
		if valid {
			r.globalRoot.insert(file)
		} else {
			duplicates = append(duplicates, node.Files[i])
		}
	}
	return duplicates, nil
}

func createFile(path string) {

}
