package mcp

import (
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"mcp-k8s/Model"
	"mcp-k8s/Repository/Cluster"
	"mcp-k8s/Util"
	"strconv"
)

func GetPods(filepath, namespace string) []Model.Pod {
	cluster, _ := Cluster.Build(filepath)
	pods, err := cluster.Pods(namespace)
	if err != nil {
		Util.Logger.Error("获取POD列表失败", zap.Error(err))
		return []Model.Pod{}
	}

	res := make([]Model.Pod, 0)

	for _, v := range pods.Items {
		pod := Model.Pod{
			PodID:          string(v.UID),
			PodName:        v.Name,
			PodIP:          v.Status.PodIP,
			Namespace:      v.Namespace,
			NodeID:         v.Spec.NodeName,
			NodeIP:         v.Status.HostIP,
			ContainersName: Util.PluckAndJoin(v.Spec.Containers, "Name", ","),
			ImagesVersion:  Util.PluckAndJoin(v.Spec.Containers, "Image", ","),
			CreateTime:     v.CreationTimestamp.Time,
			Status:         string(v.Status.Phase),
		}

		for _, container := range v.Spec.Containers {
			cpuLimit := int64(0)
			memLimit := int64(0)

			if lim, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
				cpuLimit += lim.Value()
			}
			if lim, ok := container.Resources.Limits[corev1.ResourceMemory]; ok {
				memLimit += lim.Value()
			}
			pod.CPU = strconv.Itoa(int(cpuLimit))
			pod.Memory = strconv.Itoa(int(memLimit))
		}

		if len(v.OwnerReferences) != 0 {
			pod.WorkloadID = string(v.OwnerReferences[0].UID)
			pod.WorkloadName = v.OwnerReferences[0].Name
			pod.Kind = v.OwnerReferences[0].Kind
		}

		res = append(res, pod)
	}

	return res
}
