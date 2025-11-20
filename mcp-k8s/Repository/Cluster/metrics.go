package Cluster

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"mcp-k8s/Util"
)

func Metrics(filepath, namespace string) metricsv1beta1.PodMetricsList {
	config, err := clientcmd.BuildConfigFromFlags("", filepath)
	if err != nil {
		panic(fmt.Sprintf("加载K8s配置失败: %v", err))
	}

	metricsClientset, err := versioned.NewForConfig(config)
	if err != nil {
		panic(fmt.Sprintf("创建Metrics客户端失败: %v", err))
	}

	podMetricsList, err := metricsClientset.MetricsV1beta1().PodMetricses(namespace).List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: Util.Int64Ptr(30),
	})
	if err != nil {
		panic(fmt.Sprintf("获取Pod Metrics失败: %v", err))
	}

	return *podMetricsList
}
