package common

import (
	"fmt"
)

//StorageNode is the metadata for a storage node
type StorageNode struct {
	StorageIP   string
	ClientPort  int
	CommandPort int
	Files       []string
}

//GetIndexKey give the unique key to index the node
func (s *StorageNode) GetIndexKey() string {
	return fmt.Sprintf("%s:%d/%d", s.StorageIP, s.ClientPort, s.CommandPort)
}

//GetFileTokens split file names into tokens
func (s *StorageNode) GetFileTokens() [][]string {
	res := [][]string{}
	for _, file := range s.Files {
		tokens := Tokenize(file)
		res = append(res, tokens)
	}
	return res
}
