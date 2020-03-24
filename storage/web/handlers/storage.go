package handlers

import (
	"fmt"

	"cmu.edu/dfs/storage/core"
	"github.com/gin-gonic/gin"
)

//HandleSize handle the request of returning the length of a file in bytes
func HandleSize(c *gin.Context) {
	params := &pathParams{}

	c.BindJSON(&params)
	s := core.GetStorageNode()
	if !s.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("empty path"))
		return
	}
	size, err := s.GetFileSize(params.Path)

	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"size": size})
}

//HandleRead handle the request of reading a file in byte format
func HandleRead(c *gin.Context) {
	params := &fileParams{}
	c.BindJSON(&params)
	s := core.GetStorageNode()
	if !s.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("empty path"))
		return
	}
	if !s.WithinBounds(params.Path, params.Offset, params.Length) {
		c.JSON(indexOutOfBoundsException(fmt.Sprintf("out of bounds:%d,%d",
			params.Offset, params.Length)))
		return
	}
	data, err := s.Read(params.Path, params.Offset, params.Length)

	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"data": data})
}

//HandleWrite handle the request of writeing data into a file
func HandleWrite(c *gin.Context) {
	params := &fileParams{}
	c.BindJSON(&params)
	s := core.GetStorageNode()
	if !s.ValidatePath(params.Path) {
		c.JSON(illegalArgumentError("empty path"))
		return
	}
	if params.Offset < 0 || params.Length < 0 {
		c.JSON(indexOutOfBoundsException("out of bounds"))
		return
	}
	data, err := s.Write(params.Path, params.Offset, params.Data)
	if err != nil {
		c.JSON(fileNotFoundError(err.Error()))
		return
	}
	c.JSON(200, gin.H{"success": data})
}
