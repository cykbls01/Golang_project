package Cluster

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"mcp-k8s/Util"
)

func (cluster *Cluster) DeploymentByLabel(namespace, label string) (v1.DeploymentList, error) {
	deployments, err := cluster.client.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(30),
		LabelSelector:  label,
	})
	if err != nil {
		return v1.DeploymentList{}, err
	}

	return *deployments, nil
}
