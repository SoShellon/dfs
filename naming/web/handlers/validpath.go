package handlers

import (
	"cmu.edu/dfs/naming/core"
	"github.com/gin-gonic/gin"
)

//HandleIsValidPath decide whether is a path valid
func HandleIsValidPath(c *gin.Context) {
	params := struct {
		Path string `json:"path"`
	}{}
	c.BindJSON(&params)
	r := core.GetRegistrar()
	exists := r.Exists(params.Path)
	c.JSON(200, gin.H{"success": exists})
}
