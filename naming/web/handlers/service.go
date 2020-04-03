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

func illegalArgumentError(msg string) (int, map[string]string) {
	return 404, map[string]string{
		"exception_type": "IllegalArgumentException",
		"exception_info": msg,
	}
}

type pathParams struct {
	Path string `json:"path"`
}

//HandleIsValidPath decide whether is a path valid
func HandleIsValidPath(c *gin.Context) {
	params := &pathParams{}
	c.BindJSON(&params)
	r := core.GetRegistrar()
	exists := r.Exists(params.Path)
	c.JSON(200, gin.H{"success": exists})
}

//HandleGetStorage handle the request of storage server that hosts the file
func HandleGetStorage(c *gin.Context) {
	params := &pathParams{}
	c.BindJSON(&params)
	r := core.GetRegistrar()
	if !r.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("Wrong path:" + params.Path))
		return
	}
	dir, err := r.IsDir(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	if dir {
		c.JSON(fileNotFoundError(params.Path + " is a dir"))
		return
	}
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

//HandleList handles the list
func HandleList(c *gin.Context) {
	params := &pathParams{}
	c.BindJSON(&params)
	r := core.GetRegistrar()
	if !r.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("Wrong path:" + params.Path))
		return
	}
	dir, err := r.IsDir(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	if !dir {
		c.JSON(fileNotFoundError(params.Path + " is a file"))
		return
	}
	files, err := r.ListFiles(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"files": files})
}

//HandleDelete handles the request of deleting a file or directory
func HandleDelete(c *gin.Context) {
	params := &pathParams{}
	c.BindJSON(params)
	r := core.GetRegistrar()
	if !r.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("empty path"))
		return
	}
	success, err := r.Delete(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"success": success})
}

//HandleIsDir handles the requesting of determining whether a path
//refers to a path
func HandleIsDir(c *gin.Context) {
	params := &pathParams{}
	c.BindJSON(&params)
	r := core.GetRegistrar()
	if !r.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("Wrong path:" + params.Path))
		return
	}
	isDir, err := r.IsDir(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"success": isDir})
}

//HandleCreateDirectory handle the request of creating a directory
func HandleCreateDirectory(c *gin.Context) {
	params := &pathParams{}
	c.BindJSON(&params)
	r := core.GetRegistrar()
	if !r.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("Wrong path:" + params.Path))
		return
	}
	success, err := r.CreateDir(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"success": success})
}

//HandleCreateFile handle the requesting of creating a file
func HandleCreateFile(c *gin.Context) {
	params := &pathParams{}
	c.BindJSON(&params)
	r := core.GetRegistrar()
	if !r.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("Wrong path:" + params.Path))
		return
	}
	success, err := r.CreateFile(params.Path)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"success": success})
}
