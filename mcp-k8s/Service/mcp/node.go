package mcp

import (
	"go.uber.org/zap"
	"mcp-k8s/Model"
	"mcp-k8s/Repository/Cluster"
	"mcp-k8s/Util"
	"strings"
)

func GetNodes(filepath string) []Model.Node {
	cluster, _ := Cluster.Build(filepath)
	cluster, err := Cluster.Build(filepath)
	if err != nil {
		Util.Logger.Error("构建集群失败", zap.Error(err))
		return []Model.Node{}
	}

	nodes, err := cluster.Nodes()
	res := make([]Model.Node, 0)
	if err != nil {
		if strings.Contains(err.Error(), "unexpected EOF") {
			pods, err := cluster.PodByLabel("kube-system", "app=icagent")
			if err != nil {
				Util.Logger.Error("获取Icagent列表失败", zap.Error(err))
				return []Model.Node{}
			}
			for _, pod := range pods.Items {
				res = append(res, Model.Node{
					NodeName: pod.Status.PodIP,
					NodeIP:   pod.Status.PodIP,
				})
			}
			return res
		}
		Util.Logger.Error("获取Node列表失败", zap.Error(err))
		return []Model.Node{}
	}
	for _, v := range nodes.Items {
		res = append(res, Model.Node{
			NodeName: v.Name,
			NodeIP:   v.Status.Addresses[0].Address,
		})
	}
	return res
}
