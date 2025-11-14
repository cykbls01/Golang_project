package Cluster

import (
	"context"
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"time"
)

func GetNamespaces(kubeconfigPath string, excludeNamespaces []string) ([]string, error) {
	// -------------- 第一步：校验kubeconfig有效性 --------------
	// 1. 校验kubeconfig路径不为空
	if kubeconfigPath == "" {
		return nil, errors.New("kubeconfig路径不能为空")
	}

	// 2. 校验文件存在且为普通文件
	fileInfo, err := os.Stat(kubeconfigPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("kubeconfig文件不存在: %s", kubeconfigPath)
		}
		return nil, fmt.Errorf("获取kubeconfig文件信息失败: %w", err)
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("kubeconfig路径是目录，需传入文件路径: %s", kubeconfigPath)
	}

	// 3. 校验文件可读（尝试只读打开）
	file, err := os.OpenFile(kubeconfigPath, os.O_RDONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return nil, fmt.Errorf("kubeconfig文件无读取权限: %s", kubeconfigPath)
		}
		return nil, fmt.Errorf("打开kubeconfig文件失败: %w", err)
	}
	defer file.Close()

	// 4. 校验kubeconfig格式有效（加载配置）
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("kubeconfig格式无效或配置错误: %w", err)
	}

	// -------------- 第二步：创建K8s客户端 --------------
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建K8s客户端失败: %w", err)
	}

	// -------------- 第三步：获取所有命名空间 --------------
	// 设置10秒超时上下文，避免API调用阻塞
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	namespaceList, err := clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取命名空间列表失败: %w", err)
	}

	// -------------- 第四步：过滤指定命名空间 --------------
	// 将排除列表转为map，优化查询效率（O(1)查询）
	excludeMap := make(map[string]struct{}, len(excludeNamespaces))
	for _, ns := range excludeNamespaces {
		excludeMap[ns] = struct{}{}
	}

	// 收集过滤后的命名空间
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
