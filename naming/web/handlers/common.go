package handlers

import (
	"cmu.edu/dfs/naming/core"
	"github.com/gin-gonic/gin"
)

func fileNotFoundError(msg string) (int, map[string]string) {
	return 404, map[string]string{
		"exception_type": "FileNotFoundException",
		"exception_info": msg,
	}
}

type pathParams struct {
	Path string `json:"path"`
}

//HandleList handles the list
func HandleList(c *gin.Context) {
	params := &pathParams{}
	c.BindJSON(&params)
	r := core.GetRegistrar()
	files, err := r.ListFile(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"files": files})
}

//HandleIsDir handles the requesting of determining whether a path
//refers to a path
func HandleIsDir(c *gin.Context) {
	params := &pathParams{}
	c.BindJSON(&params)
	r := core.GetRegistrar()
	isDir, err := r.IsDir(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"success": isDir})
}
