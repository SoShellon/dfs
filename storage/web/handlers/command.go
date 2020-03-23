package handlers

import (
	"cmu.edu/dfs/storage/core"
	"github.com/gin-gonic/gin"
)

//HandleCreate handle the request of creating a file from client
func HandleCreate(c *gin.Context) {
	req := &pathParams{}
	c.Bind(req)
	s := core.GetStorageNode()

	success, err := s.CreateFile(req.Path)

	if err != nil {
		c.JSON(illegalArgumentError(err.Error()))
		return
	}

	c.JSON(200, gin.H{"success": success})
}

//HandleDelete handles the request of deleting a file
func HandleDelete(c *gin.Context) {
	req := &pathParams{}
	c.Bind(req)
	s := core.GetStorageNode()

	success, err := s.DeleteFile(req.Path)

	if err != nil {
		c.JSON(illegalArgumentError(err.Error()))
		return
	}

	c.JSON(200, gin.H{"success": success})
}
