package Cluster

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"mcp-k8s/Util"
)

func (cluster *Cluster) Pods(namespace string) (v1.PodList, error) {
	pods, err := cluster.client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(30),
	})
	if err != nil {
		return v1.PodList{}, err
	}

	return *pods, nil
}

func (cluster *Cluster) PodByLabel(namespace, label string) (v1.PodList, error) {
	pods, err := cluster.client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(30),
		LabelSelector:  label,
	})
	if err != nil {
		return v1.PodList{}, err
	}

	return *pods, nil
}
