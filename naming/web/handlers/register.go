package handlers

import (
	"cmu.edu/dfs/naming/core"
	"github.com/gin-gonic/gin"
)

func registerError(msg string) (int, map[string]string) {
	return 409, map[string]string{
		"exception_info": msg,
		"exception_type": "IllegalStateException",
	}
}

//HandleRegister handle the request of
//Registering a storage server with the naming server
func HandleRegister(c *gin.Context) {
	node := &core.StorageNode{}
	err := c.Bind(node)
	if err != nil {
		c.JSON(registerError(err.Error()))
		return
	}
	registrar := core.GetRegistrar()
	duplicateFiles, err := registrar.AddStorageNode(node)
	if err != nil {
		c.JSON(registerError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"files": duplicateFiles})
}
