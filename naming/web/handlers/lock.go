package handlers

import (
	"fmt"

	"cmu.edu/dfs/naming/core"
	"github.com/gin-gonic/gin"
)

type lockParams struct {
	Path      string `json:"path"`
	Exclusive bool   `json:"exclusive"`
}

//HandleLock handles the request of locking a file or directory
func HandleLock(c *gin.Context) {
	params := &lockParams{}
	c.BindJSON(params)
	r := core.GetRegistrar()
	if !r.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("empty path"))
		return
	}
	_, err := r.IsDir(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(fmt.Sprintf("%s does not exists", params.Path)))
		return
	}
	r.Lock(params.Path, params.Exclusive)
	c.String(200, "")
}

//HandleUnlock handles the request of Unlocking a file or directory
func HandleUnlock(c *gin.Context) {
	params := &lockParams{}
	c.BindJSON(params)
	r := core.GetRegistrar()
	if !r.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("empty path"))
		return
	}
	_, err := r.IsDir(params.Path)
	if err != nil {
		c.JSON(illegalArgumentError(fmt.Sprintf("%s does not exists", params.Path)))
		return
	}
	r.Unlock(params.Path, params.Exclusive)
	c.String(200, "")
}
