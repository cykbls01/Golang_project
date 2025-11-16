package main

import (
	"fmt"
	"mcp-k8s/Service"
	"mcp-k8s/Util"
	"os"
)

func main() {
	fmt.Println(Service.McpGetEndpoints(os.Args[1]))
	Util.Output(Service.McpGetEndpoints(os.Args[1]), "endpoint.xlsx")
}
