package task

import (
	"mcp-k8s/Service/mcp"
	"mcp-k8s/Util"
)

func UpdateServices(dirpath string) {
	services := mcp.McpGetEndpoints(dirpath)
	for _, service := range services {
		Util.DB.Save(&service)
	}
}
