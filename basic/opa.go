package main

import (
	"basic/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	_ "github.com/kr/pretty"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

type Data struct {
	ClusterName string   `json:"cluster_name"`
	Namespace   string   `json:"namespace"`
	PodName     string   `json:"pod_name"`
	Privileged  bool     `json:"privileged"`
	HostPaths   []string `json:"host_paths"`
	Sysctls     []string `json:"sysctls"`
}

type Namespace struct {
	ClusterName string `json:"cluster_name"`
	ConfigPath  string `json:"config_path"`
	Namespace   string `json:"namespace"`
}

func ExecuteComplexCmd(cmdStr string) (string, string, error) {
	// 获取当前工作目录[4](@ref)
	//fmt.Println(cmdStr)
	currentDir, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("获取目录失败: %v", err)
	}

	// 创建带bash解析的命令对象[1,6](@ref)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	cmd.Dir = currentDir // 显式设置工作目录

	// 分离标准输出和错误缓冲区[3](@ref)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// 执行并等待命令完成[6](@ref)
	err = cmd.Run()
	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()

	return stdout, stderr, err
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
		if isExcluded(namespace.Name) {
			continue
		}
		res = append(res, Namespace{
			ClusterName: filepath.Base(kubeconfigPath),
			ConfigPath:  kubeconfigPath,
			Namespace:   namespace.Name,
		})
	}
	return res
}

func WriteSliceToFile(data interface{}, filename string) error {
	// 类型校验：必须为切片类型
	if reflect.ValueOf(data).Kind() != reflect.Slice {
		return os.ErrInvalid
	}

	// 序列化为带缩进的JSON格式
	formattedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// 原子化写入文件（覆盖模式+自动创建）
	return os.WriteFile(filename, formattedData, 0644)
}

func isExcluded(ns string) bool {
	excludeList := []string{"kube-system", "arms-prom", "falco", "kube-node-lease", "kube-public", "istio-system"}
	for _, excluded := range excludeList {
		if excluded == ns {
			return true
		}
	}
	return false
}

func image_check(file, base string) {
	for _, ns := range ListAllNamespaces(file) {
		os.Mkdir(base+"/"+filepath.Base(file), 0755)
		os.Mkdir(base+"/"+filepath.Base(file)+"/"+ns.Namespace, 0755)
		nsbase := base + "/" + filepath.Base(file) + "/" + ns.Namespace
		for _, image := range ListAllImages(ListAllPods(ns)) {
			cmd1 := "trivy image --skip-db-update --skip-java-db-update  " + image + " --severity=CRITICAL  --report=summary --insecure --format template --template \"@/data/.cache/trivy/html.tpl\" --output " + nsbase + "/" + filepath.Base(image) + ".html --timeout 3m -q --cache-dir /data/.cache/trivy"
			cmd2 := "trivy image --skip-db-update --skip-java-db-update  " + image + " --severity=CRITICAL  --report=summary --insecure --format json --output " + nsbase + "/" + filepath.Base(image) + ".json --timeout 3m -q --cache-dir /data/.cache/trivy"
			switch {
			case strings.HasPrefix(image, "cr-ee"):
				ExecuteComplexCmd(cmd1)
				ExecuteComplexCmd(cmd2)
			case strings.HasPrefix(image, "swr.bj-kjy-50"):
				ExecuteComplexCmd(cmd1 + " --username " + util.Config.MP["swr5-username"] + " --password " + util.Config.MP["swr-password"])
				ExecuteComplexCmd(cmd2 + " --username " + util.Config.MP["swr5-username"] + " --password " + util.Config.MP["swr-password"])
			case strings.HasPrefix(image, "swr.bj-kjy-90"):
				ExecuteComplexCmd(cmd1 + " --username " + util.Config.MP["swr9-username"] + " --password " + util.Config.MP["swr-password"])
				ExecuteComplexCmd(cmd2 + " --username " + util.Config.MP["swr9-username"] + " --password " + util.Config.MP["swr-password"])
			case strings.HasPrefix(image, "swr.sh-ky-70"):
				ExecuteComplexCmd(cmd1 + " --username " + util.Config.MP["swr7-username"] + " --password " + util.Config.MP["swr-password"])
				ExecuteComplexCmd(cmd2 + " --username " + util.Config.MP["swr7-username"] + " --password " + util.Config.MP["swr-password"])
			default:
				fmt.Println("未知前缀: " + image)
			}

		}
	}

}

func isPrivilieged(pod v1.Pod) bool {
	for _, container := range pod.Spec.Containers {
		if container.SecurityContext != nil && container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged {
			return true
		}
	}
	return false
}

func pod_check(path, base string) {
	var result = make(map[string]Data)
	namespaces := ListAllNamespaces(path)
	for _, ns := range namespaces {
		result = make(map[string]Data)
		nsName := ns.Namespace
		// 检查是否在排除列表中
		if isExcluded(nsName) {
			continue
		}

		// 获取命名空间中的所有Pod
		pods := ListAllPods(ns)

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
					if volume.HostPath.Path == "/etc/localtime" || strings.Contains(volume.HostPath.Path, "/data") || strings.Contains(volume.HostPath.Path, "/home") {
						continue
					}
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

		res := make([]Data, 0)
		for _, value := range result {
			count.HostPaths += len(value.HostPaths)
			count.Sysctls += len(value.Sysctls)
			count.Privileged += boolToInt(value.Privileged)
			res = append(res, value)
		}
		if len(res) == 0 {
			continue
		}
		os.Mkdir(base+"/"+filepath.Base(path), 0755)
		WriteSliceToFile(res, base+"/"+filepath.Base(path)+"/"+ns.Namespace+".opa")
	}
}

type Count struct {
	Privileged int
	HostPaths  int
	Sysctls    int
}

var count Count

func main() {
	util.Init()
	files, _ := ListAllFiles(util.Config.MP["path"])
	pretty.Println(files)
	os.Mkdir("pod_check_"+time.Now().Format("2006-01-02"), 0755)
	os.Mkdir("image_check_"+time.Now().Format("2006-01-02"), 0755)
	for _, file := range files {
		image_check(file, "image_check_"+time.Now().Format("2006-01-02"))
		pod_check(file, "pod_check_"+time.Now().Format("2006-01-02"))
	}
	pretty.Println(count)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
