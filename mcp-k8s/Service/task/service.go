package task

import (
	"go.uber.org/zap"
	"mcp-k8s/Service/mcp"
	"mcp-k8s/Util"
)

func UpdateServices(dirPath string) {
	files, _ := Util.WalkAllFilesAbs(dirPath)
	for _, file := range files {
		Util.Logger.Info("process", zap.String("file", file))
		Util.SaveFromFunc(file, mcp.GetEndpoints)
	}
}
