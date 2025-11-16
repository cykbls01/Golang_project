package Cluster

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log/slog"
	"mcp-k8s/Util"
)

func Namespaces(filepath string, excludeNamespaces []string) ([]string, error) {
	coreClientset, err := Build(filepath)
	if err != nil {
		slog.Error("创建核心客户端失败: %v", err)
		return []string{}, err
	}

	namespaceList, err := coreClientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(30),
	})
	if err != nil {
		slog.Error("获取Namespace列表失败: %v", err)
		return []string{}, err
	}

	excludeMap := make(map[string]struct{}, len(excludeNamespaces))
	for _, ns := range excludeNamespaces {
		excludeMap[ns] = struct{}{}
	}

	var result []string
	for _, ns := range namespaceList.Items {
		nsName := ns.Name
		// 跳过排除列表中的命名空间
		if _, excluded := excludeMap[nsName]; !excluded {
			result = append(result, nsName)
		}
	}
	return result, nil
}
