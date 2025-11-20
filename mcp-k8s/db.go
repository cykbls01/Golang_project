package main

import (
	"mcp-k8s/Service/task"
	"mcp-k8s/Util"
	"os"
)

func main() {
	Util.Pre()
	Util.InitDB()
	Util.DB.Exec("Truncate Table service")
	task.UpdateServices(os.Args[1])
	Util.Post()
}
