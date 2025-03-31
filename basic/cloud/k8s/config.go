package k8s

import (
	"basic/cloud/hcs"
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"time"
)

func GenerateConfigs() {
	for _, p := range hcs.ListProject() {
		for _, c := range hcs.ListCluster(p.Id) {
			config := hcs.GetKubeConfig(p.Id, *c.Metadata.Uid)
			jsonData, err := json.Marshal(config)
			if err != nil {
				fmt.Println("JSON marshaling failed:", err)
				return
			}

			// 将JSON数据写入文件
			err = os.WriteFile(c.Metadata.Name+".json", jsonData, 0644)
			if err != nil {
				fmt.Println("Failed to write JSON data to file:", err)
				return
			}
			break
		}
	}
}

func ValidateKubeconfig(kubeconfigPath string) error {
	// 第一阶段：基础文件验证
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		return fmt.Errorf("kubeconfig 文件不存在: %v", err)
	}

	// 第二阶段：配置加载验证（含证书解析）[3,6](@ref)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return fmt.Errorf("配置文件解析失败: %v", err)
	}

	// 第三阶段：API Server 连通性验证
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("创建客户端失败: %v", err)
	}

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 执行简单 API 请求验证凭证有效性[3](@ref)
	if _, err := clientset.Discovery().ServerVersion(); err != nil {
		return fmt.Errorf("API Server 连接失败: %v", err)
	}

	return nil
}
