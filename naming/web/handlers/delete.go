package handlers

import (
	"cmu.edu/dfs/naming/core"
	"github.com/gin-gonic/gin"
)

//HandleDelete handles the request of deleting a file or directory
func HandleDelete(c *gin.Context) {
	params := struct {
		Path string `json:"Path"`
	}{}
	r := core.GetRegistrar()
	success, err := r.Delete(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"success": success})
}
