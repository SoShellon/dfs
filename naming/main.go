package main

import (
	"fmt"
	"os"
	"strconv"

	"cmu.edu/dfs/common"
	"cmu.edu/dfs/naming/core"
	"cmu.edu/dfs/naming/web"
)

const defaultListenAddress = "localhost:%d"

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		common.Error(fmt.Sprintf("The length of arguments should be 2:%+v", args))
		os.Exit(-1)
		// args = []string{"8080", "8081"}
	}
	servicePort, err1 := strconv.Atoi(args[0])
	registrationPort, err2 := strconv.Atoi(args[1])

	if err1 != nil || err2 != nil {
		common.Error(fmt.Sprintf("Wrong arguments:%s", args))
		os.Exit(-1)
	}
	core.InitRegistrar()
	server := web.GetNamingServer(servicePort, registrationPort)
	server.Run()
}
