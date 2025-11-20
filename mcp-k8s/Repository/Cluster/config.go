package Cluster

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"mcp-k8s/Util"
	path "path/filepath"
)

type Cluster struct {
	client kubernetes.Clientset
}

func Build(filepath string) (Cluster, error) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath)
	if err != nil {
		return Cluster{}, err
	}

	coreClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return Cluster{}, err
	}

	_, err = coreClientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(5),
	})
	if err != nil {
		return Cluster{}, err
	}

	return Cluster{client: *coreClientset}, nil
}

func GetName(filepath string) string {
	return path.Base(filepath)
}
