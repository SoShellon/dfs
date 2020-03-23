package main

import (
	"os"
	"strconv"

	"cmu.edu/dfs/common"
	"cmu.edu/dfs/storage/core"
	"cmu.edu/dfs/storage/web"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		common.Error("The length of arguments should be 3")
		// os.Exit(-1)
		args = []string{"8082", "8083", "8081", "/tmp/xianglol"}
	}
	clientPort, _ := strconv.Atoi(args[0])
	commandPort, _ := strconv.Atoi(args[1])
	namingPort, _ := strconv.Atoi(args[2])
	dirPath := args[3]
	core.InitStorageNode(commandPort, namingPort, namingPort, dirPath)
	server := web.GetStorageServer(clientPort, commandPort)
	server.Run()
}
