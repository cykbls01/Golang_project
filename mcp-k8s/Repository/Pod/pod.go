package Pod

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log/slog"
	"mcp-k8s/Repository/Cluster"
	"mcp-k8s/Util"
)

func Pods(filepath, namespace string) (v1.PodList, error) {

	coreClientset, err := Cluster.Build(filepath)
	if err != nil {
		slog.Error("创建核心客户端失败: %v", err)
		return v1.PodList{}, err
	}

	pods, err := coreClientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(30),
	})
	if err != nil {
		slog.Error("获取Pod列表失败: %v", err)
		return v1.PodList{}, err
	}

	return *pods, nil
}

func PodByLabel(filepath, namespace, label string) (v1.PodList, error) {

	coreClientset, err := Cluster.Build(filepath)
	if err != nil {
		slog.Error("创建核心客户端失败: %v", err)
		return v1.PodList{}, err
	}

	pods, err := coreClientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(30),
		LabelSelector:  label,
	})
	if err != nil {
		slog.Error("获取Pod列表失败: %v", err)
		return v1.PodList{}, err
	}

	return *pods, nil
}
