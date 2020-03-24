package web

import (
	"cmu.edu/dfs/common"
	"cmu.edu/dfs/storage/web/handlers"
	"github.com/gin-gonic/gin"
)

type storageServer struct {
	clientPort  int
	commandPort int
}

//GetStorageServer build a new storage server instance
func GetStorageServer(clientPort int, commandPort int) common.Server {
	return &storageServer{
		clientPort:  clientPort,
		commandPort: commandPort,
	}
}
func (s *storageServer) runClientServer() error {
	r := gin.Default()
	r.POST("/storage_size", handlers.HandleSize)
	r.POST("/storage_read", handlers.HandleRead)
	r.POST("/storage_write", handlers.HandleWrite)
	return common.RunServer(r, s.clientPort)
}

func (s *storageServer) runCommandServer() error {
	r := gin.Default()
	r.POST("/storage_create", handlers.HandleCreate)
	r.POST("/storage_delete", handlers.HandleDelete)
	r.POST("/storage_copy", handlers.HandleCopy)
	return common.RunServer(r, s.commandPort)
}
func (s *storageServer) Run() {
	joiner := make(chan error)
	go func() { joiner <- s.runCommandServer() }()
	go func() { joiner <- s.runClientServer() }()
	err := <-joiner
	if err != nil {
		common.Error(err.Error())
	}
}
