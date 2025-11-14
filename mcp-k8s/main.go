package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"mcp-k8s/Resources/Pod"
	"os"
)

// PodResourceUsage 存储Pod资源使用信息
type PodResourceUsage struct {
	Namespace    string
	PodName      string
	CPURequest   string  // 请求CPU（格式化后）
	CPULimit     string  // 限制CPU（格式化后）
	CPUUsage     string  // 实际使用CPU（格式化后）
	CPUUsageRate float64 // CPU使用率（相对于限制，百分比）
	MemRequest   string  // 请求内存（格式化后）
	MemLimit     string  // 限制内存（格式化后）
	MemUsage     string  // 实际使用内存（格式化后）
	MemUsageRate float64 // 内存使用率（相对于限制，百分比）
	PodStatus    corev1.PodPhase
	StartTime    metav1.Time
}

func main() {
	kubeconfig := clientcmd.RecommendedHomeFile
	if envvar := os.Getenv("KUBECONFIG"); envvar != "" {
		kubeconfig = envvar
	}
	Pod.Get(kubeconfig, "cram")
}
