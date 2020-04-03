package web

import (
	"cmu.edu/dfs/common"
	"cmu.edu/dfs/naming/web/handlers"
	"github.com/gin-gonic/gin"
)

type namingServer struct {
	servicePort       int
	registerationPort int
}

//GetNamingServer help client to get a naming server instance
func GetNamingServer(servicePort int, registerationPort int) common.Server {
	return &namingServer{servicePort, registerationPort}
}

func (n *namingServer) runService() error {
	engine := gin.Default()
	engine.POST("/is_valid_path", handlers.HandleIsValidPath)
	engine.POST("/getstorage", handlers.HandleGetStorage)
	engine.POST("/delete", handlers.HandleDelete)
	engine.POST("/create_directory", handlers.HandleCreateDirectory)
	engine.POST("/create_file", handlers.HandleCreateFile)
	engine.POST("/list", handlers.HandleList)
	engine.POST("/is_directory", handlers.HandleIsDir)
	engine.POST("/unlock", handlers.HandleUnlock)
	engine.POST("/lock", handlers.HandleLock)
	return common.RunServer(engine, n.servicePort)
}

func (n *namingServer) runRegistration() error {
	engine := gin.Default()
	engine.POST("/register", handlers.HandleRegister)
	return common.RunServer(engine, n.registerationPort)
}
func (n *namingServer) Run() {
	joiner := make(chan error)
	go func() { joiner <- n.runService() }()
	go func() { joiner <- n.runRegistration() }()
	err := <-joiner
	common.Error(err.Error())
}
