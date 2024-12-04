package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getNamespacesFromK8sCluster(kubeconfigPath string) ([]string, error) {
	config, err := getKubeConfig(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %v", err)
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %v", err)
	}

	var namespaceList []string
	for _, ns := range namespaces.Items {
		namespaceList = append(namespaceList, ns.Name)
	}

	return namespaceList, nil
}

func getKubeConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		kubeconfigPath = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err == nil {
		return config, nil
	}
	return nil, err
}

func main() {
	// 指定每个集群的kubeconfig文件路径
	kubeconfigPaths := []string{"kubeconfig1", "kubeconfig2"}

	// 创建一个map来保存命名空间数据
	namespacesMap := make(map[string][]string)

	// 遍历集群配置，获取每个集群的命名空间
	for _, kubeconfigPath := range kubeconfigPaths {
		namespaces, err := getNamespacesFromK8sCluster(kubeconfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting namespaces from cluster: %v\n", err)
			os.Exit(1)
		}
		namespacesMap[kubeconfigPath] = namespaces
	}

	// 打印结果
	for kubeconfigPath, namespaces := range namespacesMap {
		fmt.Printf("Namespaces from cluster with kubeconfig '%s': %v\n", kubeconfigPath, namespaces)
	}
}
