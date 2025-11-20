package mcp

import (
	"go.uber.org/zap"
	"mcp-k8s/Model"
	"mcp-k8s/Repository/Cluster"
	"mcp-k8s/Util"
)

func GetImages(filepath, namespace string) []Model.Image {
	cluster, _ := Cluster.Build(filepath)
	pods, err := cluster.Pods(namespace)
	if err != nil {
		Util.Logger.Error("获取POD列表失败", zap.Error(err))
		return []Model.Image{}
	}

	imageList := make([]Model.Image, 0)
	for _, pod := range pods.Items {
		for _, container := range append(pod.Spec.Containers, pod.Spec.InitContainers...) {
			workloadName := ""
			workloadType := ""
			if len(pod.OwnerReferences) != 0 {
				workloadName = pod.OwnerReferences[0].Name
				workloadType = pod.OwnerReferences[0].Kind
			}

			imageList = append(imageList, Model.Image{
				Namespace:    pod.Namespace,
				PodName:      pod.Name,
				ImageName:    container.Image,
				WorkloadName: workloadName,
				WorkloadType: workloadType,
				Cluster:      Cluster.GetName(filepath),
			})
		}
	}

	return imageList
}

func McpGetImages(dirPath string) []Model.Image {
	files, _ := Util.WalkAllFilesAbs(dirPath)
	res := make([]Model.Image, 0)
	for _, file := range files {
		res = append(res, GetImages(file, "")...)
	}
	return res
}
