package main

import (
	"context"
	"fmt"
	"github.com/kr/pretty"
	_ "github.com/kr/pretty"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

type Data struct {
	ClusterName string   `json:"cluster_name"`
	Namespace   string   `json:"namespace"`
	PodName     string   `json:"pod_name"`
	Privileged  bool     `json:"privileged"`
	HostPaths   []string `json:"host_paths"`
}

func ListAllFiles(dirPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录，只添加文件
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

var result map[string]Data

func main() {
	files, _ := ListAllFiles(os.Args[1])
	result = make(map[string]Data)
	for _, file := range files {
		Run(file)
	}
	pretty.Println(result)
}

func Run(kubeconfigPath string) {

	excludeNamespaces := []string{"kube-system", "arms-prom", "falco"} // 默认排除的系统命名空间

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 获取所有命名空间
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 遍历所有命名空间
	for _, ns := range namespaces.Items {
		nsName := ns.Name
		// 检查是否在排除列表中
		if isExcluded(nsName, excludeNamespaces) {
			continue
		}

		// 获取命名空间中的所有Pod
		pods, err := clientset.CoreV1().Pods(nsName).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error getting pods in namespace %s: %v\n", nsName, err)
			continue
		}

		// 遍历Pod中的容器
		for _, pod := range pods.Items {
			podName := pod.Name
			//foundIssues := false

			// 检查容器安全上下文
			for _, container := range pod.Spec.Containers {
				// 检查特权模式
				if container.SecurityContext != nil && container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged {
					result[podName] = Data{
						ClusterName: filepath.Base(kubeconfigPath),
						Namespace:   nsName,
						PodName:     podName,
						Privileged:  true,
						HostPaths:   nil,
					}
				}
			}

			// 更准确地检查hostPath卷
			for _, volume := range pod.Spec.Volumes {
				if volume.HostPath != nil {
					if volume.HostPath.Path == "/etc/localtime" {
						continue
					}
					// 找到使用hostPath卷的容器
					if value, ok := result[podName]; ok {
						value.HostPaths = append(value.HostPaths, volume.HostPath.Path)
					} else {
						result[podName] = Data{
							ClusterName: filepath.Base(kubeconfigPath),
							Namespace:   nsName,
							PodName:     podName,
							Privileged:  false,
							HostPaths:   []string{volume.HostPath.Path},
						}
					}
				}
			}
		}
	}
}

// isExcluded 检查命名空间是否在排除列表中
func isExcluded(ns string, excludeList []string) bool {
	for _, excluded := range excludeList {
		if excluded == ns {
			return true
		}
	}
	return false
}
