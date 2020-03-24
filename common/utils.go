package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const defaultListenAddress = "localhost:%d"

//Error print all the error log
func Error(str string) {
	fmt.Println(str)
}

//Log print all the log strings
func Log(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// Server is the abstract interface for servers in different roles
type Server interface {
	Run()
}

// RunServer run the server on local host
func RunServer(engine *gin.Engine, port int) error {
	return engine.Run(fmt.Sprintf(defaultListenAddress, port))
}

// Tokenize split file path into sequential path token
func Tokenize(filePath string) []string {
	filePath = strings.TrimSuffix(filePath, "/")
	tokens := strings.Split(filePath, "/")
	res := []string{}
	for _, token := range tokens {
		if len(token) > 0 {
			res = append(res, token)
		}
	}
	return res
}

// SendRequest send request to other servers and receive the response
func SendRequest(url string, data interface{}, res interface{}) (err error) {

	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", "http://"+url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, res)
	return
}
