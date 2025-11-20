package Cluster

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"mcp-k8s/Util"
)

func (cluster *Cluster) Services(namespace string) (v1.ServiceList, error) {
	services, err := cluster.client.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(30),
	})
	if err != nil {
		return v1.ServiceList{}, err
	}

	return *services, nil
}
