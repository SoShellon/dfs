package handlers

import "github.com/gin-gonic/gin"

type request struct {
	Path string 
}

//HandleCreate handle the request of creating a file from client
func HandleCreate(c *gin.Context)  {
	req := &request{}
	err :=c.Bind(req)
	if err==nil {

	}
}