package Cluster

import "C"
import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"mcp-k8s/Util"
)

func (cluster *Cluster) Namespaces() (v1.NamespaceList, error) {
	namespaceList, err := cluster.client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(30),
	})
	if err != nil {
		return v1.NamespaceList{}, err
	}
	return *namespaceList, nil
}
