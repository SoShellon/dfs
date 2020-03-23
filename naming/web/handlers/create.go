package handlers

import (
	"cmu.edu/dfs/naming/core"
	"github.com/gin-gonic/gin"
)

//HandleCreateDirectory handle the request of creating a directory
func HandleCreateDirectory(c *gin.Context) {
	params := struct {
		Path string `json:"path"`
	}{}
	c.BindJSON(&params)
	r := core.GetRegistrar()
	success, err := r.CreateDir(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"success": success})
}

//HandleCreateFile handle the requesting of creating a file
func HandleCreateFile(c *gin.Context) {
	params := struct {
		Path string `json:"path"`
	}{}
	c.BindJSON(&params)
	r := core.GetRegistrar()
	success, err := r.CreateFile(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"succcess": success})
}
