package Cluster

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log/slog"
	"mcp-k8s/Util"
)

func Nodes(filepath string) (v1.NodeList, error) {
	coreClientset, err := Build(filepath)
	if err != nil {
		slog.Error("创建核心客户端失败: %v", err)
		return v1.NodeList{}, err
	}

	nodes, err := coreClientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(30),
	})
	if err != nil {
		slog.Error("获取Node列表失败: %v", err)
		return v1.NodeList{}, err
	}
	return *nodes, err
}
