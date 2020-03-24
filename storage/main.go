package main

import (
	"fmt"
	"os"
	"strconv"

	"cmu.edu/dfs/common"
	"cmu.edu/dfs/storage/core"
	"cmu.edu/dfs/storage/web"
)

func main() {
	args := os.Args[1:]
	if len(args) != 4 {
		common.Error(fmt.Sprintf("The length of arguments should be 3:%+v", args))
		os.Exit(-1)
		// args = []string{"8082", "8083", "8081", "/tmp/xianglol"}
	}
	clientPort, _ := strconv.Atoi(args[0])
	commandPort, _ := strconv.Atoi(args[1])
	namingPort, _ := strconv.Atoi(args[2])
	dirPath := args[3]
	err := core.InitStorageNode(clientPort, commandPort, namingPort, dirPath)
	if err != nil {
		common.Error(fmt.Sprintf("Storage server cannot start up:%+v", err))
		os.Exit(-1)
	}
	server := web.GetStorageServer(clientPort, commandPort)
	server.Run()
}
