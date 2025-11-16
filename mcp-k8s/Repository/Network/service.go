package Network

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log/slog"
	"mcp-k8s/Repository/Cluster"
	"mcp-k8s/Util"
)

func Services(filepath, namespace string) (v1.ServiceList, error) {
	coreClientset, err := Cluster.Build(filepath)
	if err != nil {
		slog.Error("创建核心客户端失败: %v", err)
		return v1.ServiceList{}, err
	}

	services, err := coreClientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(30),
	})
	if err != nil {
		slog.Error("获取Service列表失败: %v", err)
		return v1.ServiceList{}, err
	}

	return *services, nil
}
