package k8s

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
)

type Namespace struct {
	ClusterName string `json:"cluster_name"`
	ConfigPath  string `json:"config_path"`
	Namespace   string `json:"namespace"`
}

func ListAllPods(namespace Namespace) []v1.Pod {
	config, err := clientcmd.BuildConfigFromFlags("", namespace.ConfigPath)
	if err != nil {
		fmt.Println(err.Error())
		return []v1.Pod{}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
		return []v1.Pod{}
	}

	pods, _ := clientset.CoreV1().Pods(namespace.Namespace).List(context.TODO(), metav1.ListOptions{})
	return pods.Items
}

func ListAllNamespaces(kubeconfigPath string) []Namespace {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Println(err.Error())
		return []Namespace{}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
		return []Namespace{}
	}

	// 获取所有命名空间
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
		return []Namespace{}
	}

	var res = make([]Namespace, 0)
	for _, namespace := range namespaces.Items {
		res = append(res, Namespace{
			ClusterName: filepath.Base(kubeconfigPath),
			ConfigPath:  kubeconfigPath,
			Namespace:   namespace.Name,
		})
	}
	return res
}

func ListAllImages(pods []v1.Pod) []string {
	// 使用 map 实现去重
	imageSet := make(map[string]struct{})
	var uniqueImages []string

	// 遍历所有 Pod 的容器
	for _, pod := range pods {
		// 获取普通容器
		for _, container := range pod.Spec.Containers {
			if _, exists := imageSet[container.Image]; !exists {
				imageSet[container.Image] = struct{}{}
				uniqueImages = append(uniqueImages, container.Image)
			}
		}

		// 获取初始化容器[2](@ref)
		for _, initContainer := range pod.Spec.InitContainers {
			if _, exists := imageSet[initContainer.Image]; !exists {
				imageSet[initContainer.Image] = struct{}{}
				uniqueImages = append(uniqueImages, initContainer.Image)
			}
		}
	}

	return uniqueImages
}
