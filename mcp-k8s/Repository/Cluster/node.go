package Cluster

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"mcp-k8s/Util"
)

func (cluster *Cluster) Nodes() (v1.NodeList, error) {
	nodes, err := cluster.client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(30),
	})
	if err != nil {
		return v1.NodeList{}, err
	}
	return *nodes, err
}

func (cluster *Cluster) Node(name string) (v1.Node, error) {
	node, err := cluster.client.CoreV1().Nodes().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return v1.Node{}, err
	}
	return *node, err
}
