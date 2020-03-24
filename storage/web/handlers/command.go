package handlers

import (
	"cmu.edu/dfs/storage/core"
	"github.com/gin-gonic/gin"
)

//HandleCreate handle the request of creating a file from client
func HandleCreate(c *gin.Context) {
	params := &pathParams{}
	c.BindJSON(params)
	s := core.GetStorageNode()
	if !s.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("empty path"))
		return
	}
	success, err := s.CreateFile(params.Path)

	if err != nil {
		c.JSON(illegalArgumentError(err.Error()))
		return
	}

	c.JSON(200, gin.H{"success": success})
}

//HandleDelete handles the request of deleting a file
func HandleDelete(c *gin.Context) {
	params := &pathParams{}
	c.BindJSON(params)
	s := core.GetStorageNode()
	if !s.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("empty path"))
		return
	}
	success, err := s.DeleteFile(params.Path)

	if err != nil {
		c.JSON(illegalArgumentError(err.Error()))
		return
	}

	c.JSON(200, gin.H{"success": success})
}

//HandleCopy handle the request of copying a file from another storage server
func HandleCopy(c *gin.Context) {
	// params := &copyParams {}
	// c.BindJSON(params)
	// s := core.GetStorageNode()

}
