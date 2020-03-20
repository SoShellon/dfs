package web

import (
	"cmu.edu/dfs/common"
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
	return common.RunServer(r, s.clientPort)
}

func (s *storageServer) runCommandServer() error {
	r := gin.Default()
	// r.POST("/storage_create", handlers.HandleRegister)
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
