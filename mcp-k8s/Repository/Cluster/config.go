package Cluster

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log/slog"
)

func Build(filepath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath)
	if err != nil {
		slog.Error("加载K8s配置失败: %v", err)
		return nil, err
	}

	coreClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		slog.Error("创建核心客户端失败: %v", err)
		return nil, err
	}
	return coreClientset, nil
}
