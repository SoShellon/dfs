package handlers

import (
	"cmu.edu/dfs/naming/core"
	"github.com/gin-gonic/gin"
)

func buildErrorResponse(msg string) map[string]string {
	return map[string]string{
		"exception_type": msg,
		"exception_info": "IllegalStateException",
	}
}

//HandleRegister handle the request of
//Registering a storage server with the naming server
func HandleRegister(c *gin.Context) {
	node := &core.StorageNode{}
	err := c.Bind(node)
	if err != nil {
		c.JSON(409, buildErrorResponse(err.Error()))
		return
	}
	registrar := core.GetRegistrar()
	duplicateFiles, err := registrar.AddStorageNode(node)
	if err != nil {
		c.JSON(409, buildErrorResponse(err.Error()))
		return
	}
	c.JSON(200, gin.H{"files": duplicateFiles})
}
