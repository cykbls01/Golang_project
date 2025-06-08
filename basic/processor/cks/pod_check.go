package cks

import (
	"basic/excel"
	"basic/util/k8s"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"os"
	"path/filepath"
)

type Data struct {
	ClusterName string   `json:"cluster_name"`
	Namespace   string   `json:"namespace"`
	PodName     string   `json:"pod_name"`
	Privileged  bool     `json:"privileged"`
	HostPaths   []string `json:"host_paths"`
	Sysctls     []string `json:"sysctls"`
}

type PodCheck struct {
	Files  []string
	Path   string
	Result []Data
}

func (pc *PodCheck) Pre() {
	os.Mkdir(pc.Path, 0755)
}

func (pc *PodCheck) Process() {
	for _, file := range pc.Files {
		pc.Result = append(pc.Result, CheckPod(file, pc.Path)...)
	}
}

func (pc *PodCheck) Post() {
	fmt.Println(excel.Output(pc.Result, pc.Path+"/opa.xlsx"))
}

func isPrivilieged(pod v1.Pod) bool {
	for _, container := range pod.Spec.Containers {
		if container.SecurityContext != nil && container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged {
			return true
		}
	}
	return false
}

func CheckPod(path, base string) []Data {
	var result = make(map[string]Data)
	namespaces := k8s.ListAllNamespaces(path)
	res := make([]Data, 0)
	for _, ns := range namespaces {
		result = make(map[string]Data)
		nsName := ns.Namespace
		// 检查是否在排除列表中
		if isExcluded(nsName) {
			continue
		}

		// 获取命名空间中的所有Pod
		pods := k8s.ListAllPods(ns)

		// 遍历Pod中的容器
		for _, pod := range pods {
			podName := pod.Name
			if isPrivilieged(pod) {
				result[podName] = Data{
					ClusterName: filepath.Base(path),
					Namespace:   nsName,
					PodName:     podName,
					Privileged:  true,
					HostPaths:   make([]string, 0),
					Sysctls:     make([]string, 0),
				}
			}

			// 更准确地检查hostPath卷
			for _, volume := range pod.Spec.Volumes {
				if volume.HostPath != nil {
					//if volume.HostPath.Path == "/etc/localtime" || strings.Contains(volume.HostPath.Path, "/data") || strings.Contains(volume.HostPath.Path, "/home") {
					//	continue
					//}
					// 找到使用hostPath卷的容器
					if value, ok := result[podName]; ok {
						value.HostPaths = append(value.HostPaths, volume.HostPath.Path)
					} else {
						result[podName] = Data{
							ClusterName: filepath.Base(path),
							Namespace:   nsName,
							PodName:     podName,
							Privileged:  false,
							HostPaths:   []string{volume.HostPath.Path},
							Sysctls:     make([]string, 0),
						}
					}
				}
			}
			if pod.Spec.SecurityContext != nil && pod.Spec.SecurityContext.Sysctls != nil {
				for _, sysctl := range pod.Spec.SecurityContext.Sysctls {
					if value, ok := result[podName]; ok {
						value.Sysctls = append(value.Sysctls, sysctl.Name+":"+sysctl.Value)
					} else {
						result[podName] = Data{
							ClusterName: filepath.Base(path),
							Namespace:   nsName,
							PodName:     podName,
							Privileged:  true,
							HostPaths:   make([]string, 0),
							Sysctls:     []string{sysctl.Name + ":" + sysctl.Value},
						}
					}
				}
			}

		}
		for _, v := range result {
			res = append(res, v)
		}
	}
	return res
}
