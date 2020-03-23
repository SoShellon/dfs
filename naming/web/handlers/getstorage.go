package handlers

import (
	"cmu.edu/dfs/naming/core"
	"github.com/gin-gonic/gin"
)

//HandleGetStorage handle the request of storage server that hosts the file
func HandleGetStorage(c *gin.Context) {
	params := struct {
		Path string `json:"path"`
	}{}
	c.Bind(&params)
	r := core.GetRegistrar()
	storageNode, err := r.GetStorageNode(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
	} else {
		c.JSON(200, gin.H{
			"server_ip":   storageNode.StorageIP,
			"server_port": storageNode.ClientPort,
		})
	}
}
