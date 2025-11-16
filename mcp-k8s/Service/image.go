package Service

import (
	"mcp-k8s/Repository/Pod"
	"mcp-k8s/Util"
)

func GetImages(filepath, namespace string) []string {

	pods, _ := Pod.Pods(filepath, namespace)
	imageSet := make(map[string]struct{})

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if container.Image != "" { // 过滤空镜像（理论上不会存在）
				imageSet[container.Image] = struct{}{}
			}
		}
		for _, initContainer := range pod.Spec.InitContainers {
			if initContainer.Image != "" {
				imageSet[initContainer.Image] = struct{}{}
			}
		}
	}

	uniqueImages := make([]string, 0, len(imageSet))
	for image := range imageSet {
		uniqueImages = append(uniqueImages, image)
	}

	return uniqueImages
}

func McpGetImages(dirPath string) []string {
	files, _ := Util.WalkAllFilesAbs(dirPath)
	res := make([]string, 0)
	for _, file := range files {
		res = append(res, GetImages(file, "")...)
	}
	return res
}
