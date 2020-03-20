package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const defaultListenAddress = "localhost:%d"

//Error print all the error log
func Error(str string) {
	fmt.Println(str)
}

// Server is the abstract interface for servers in different roles
type Server interface {
	Run()
}

// RunServer run the server on local host
func RunServer(engine *gin.Engine, port int) error {
	return engine.Run(fmt.Sprintf(defaultListenAddress, port))
}
