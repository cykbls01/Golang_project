package Pod

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"math"
	"os"
	"sort"
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

func Get() {
	// 1. 加载K8s配置
	kubeconfig := clientcmd.RecommendedHomeFile
	if envvar := os.Getenv("KUBECONFIG"); envvar != "" {
		kubeconfig = envvar
	}

	// 构建K8s核心客户端配置
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(fmt.Sprintf("加载K8s配置失败: %v", err))
	}

	// 2. 创建客户端
	// 核心客户端（用于获取Pod基本信息）
	coreClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Sprintf("创建核心客户端失败: %v", err))
	}

	// Metrics客户端（用于获取Pod资源使用数据）
	metricsClientset, err := versioned.NewForConfig(config)
	if err != nil {
		panic(fmt.Sprintf("创建Metrics客户端失败: %v", err))
	}

	// 3. 获取所有Namespace的Pod
	pods, err := coreClientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: int64Ptr(30),
	})
	if err != nil {
		panic(fmt.Sprintf("获取Pod列表失败: %v", err))
	}

	// 4. 获取所有Pod的Metrics数据
	podMetricsList, err := metricsClientset.MetricsV1beta1().PodMetricses("").List(context.Background(), metav1.ListOptions{
		TimeoutSeconds: int64Ptr(30),
	})
	if err != nil {
		panic(fmt.Sprintf("获取Pod Metrics失败: %v", err))
	}

	// 构建Metrics缓存（namespace/podName -> PodMetrics）
	metricsCache := make(map[string]metricsv1beta1.PodMetrics)
	for _, pm := range podMetricsList.Items {
		key := fmt.Sprintf("%s/%s", pm.Namespace, pm.Name)
		metricsCache[key] = pm
	}

	// 5. 计算每个Pod的资源使用率
	var usageList []PodResourceUsage
	for _, pod := range pods.Items {
		// 跳过已终止的Pod
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			continue
		}

		// 从缓存获取当前Pod的Metrics
		metricsKey := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
		podMetrics, hasMetrics := metricsCache[metricsKey]

		// 计算Pod级别的资源请求/限制（汇总所有容器）
		cpuRequest := int64(0)
		cpuLimit := int64(0)
		memRequest := int64(0)
		memLimit := int64(0)

		for _, container := range pod.Spec.Containers {
			// CPU资源（单位：nano cores，1 CPU = 1e9 nano cores）
			if req, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
				cpuRequest += req.Value()
			}
			if lim, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
				cpuLimit += lim.Value()
			}

			// 内存资源（单位：bytes）
			if req, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
				memRequest += req.Value()
			}
			if lim, ok := container.Resources.Limits[corev1.ResourceMemory]; ok {
				memLimit += lim.Value()
			}
		}

		// 计算实际使用量（汇总所有容器）
		var cpuUsage int64
		var memUsage int64
		if hasMetrics {
			for _, containerMetrics := range podMetrics.Containers {
				cpuUsage += containerMetrics.Usage.Cpu().MilliValue()
				memUsage += containerMetrics.Usage.Memory().Value()
			}
		}

		// 计算使用率（相对于限制，无限制则显示"无限制"）
		cpuUsageRate := 0.0
		if cpuLimit > 0 && hasMetrics {
			cpuUsageRate = math.Round((float64(cpuUsage)/float64(cpuLimit*1000))*100*100) / 100 // 保留2位小数
		}

		memUsageRate := 0.0
		if memLimit > 0 && hasMetrics {
			memUsageRate = math.Round((float64(memUsage)/float64(memLimit))*100*100) / 100
		}

		// 格式化输出（转换为人类可读单位）
		usageList = append(usageList, PodResourceUsage{
			Namespace:    pod.Namespace,
			PodName:      pod.Name,
			CPURequest:   formatCPU(cpuRequest * 1000),
			CPULimit:     formatCPU(cpuLimit * 1000),
			CPUUsage:     formatCPU(cpuUsage),
			CPUUsageRate: cpuUsageRate,
			MemRequest:   formatMemory(memRequest),
			MemLimit:     formatMemory(memLimit),
			MemUsage:     formatMemory(memUsage),
			MemUsageRate: memUsageRate,
			PodStatus:    pod.Status.Phase,
			StartTime:    *pod.Status.StartTime,
		})
	}

	// 6. 按CPU使用率降序排序
	sort.Slice(usageList, func(i, j int) bool {
		return usageList[i].CPUUsageRate > usageList[j].CPUUsageRate
	})

	// 7. 打印结果
	printUsageTable(usageList)
}

// 格式化CPU单位（nano cores -> m/CPU）
func formatCPU(milliCores int64) string {
	//if nanoCores == 0 {
	//	return "0m"
	//}
	//milliCores := nanoCores / 1e6 // 转换为毫核（1 CPU = 1000m）
	//if milliCores >= 1000 {
	//	return fmt.Sprintf("%.1f", float64(milliCores)/1000)
	//}
	//return fmt.Sprintf("%dm", milliCores)
	return fmt.Sprintf("%dm", milliCores)
}

// 格式化内存单位（bytes -> KiB/MiB/GiB）
func formatMemory(bytes int64) string {
	if bytes == 0 {
		return "0Mi"
	}
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	unitIndex := 0
	value := float64(bytes)

	for value >= 1024 && unitIndex < len(units)-1 {
		value /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.1f%s", value, units[unitIndex])
}

// 打印资源使用表格
func printUsageTable(usageList []PodResourceUsage) {
	// 表头
	fmt.Printf("%-15s %-30s %-10s %-10s %-10s %-8s %-12s %-12s %-12s %-8s %-20s\n",
		"命名空间", "Pod名称", "CPU请求", "CPU限制", "CPU使用", "CPU使用率(%)",
		"内存请求", "内存限制", "内存使用", "内存使用率(%)", "状态")

	fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------------------------------")

	// 表内容
	for _, usage := range usageList {
		cpuRateStr := fmt.Sprintf("%.2f", usage.CPUUsageRate)
		if usage.CPULimit == "0m" {
			cpuRateStr = "无限制"
		}

		memRateStr := fmt.Sprintf("%.2f", usage.MemUsageRate)
		if usage.MemLimit == "0Mi" {
			memRateStr = "无限制"
		}

		fmt.Printf("%-15s %-30s %-10s %-10s %-10s %-8s %-12s %-12s %-12s %-8s %-20s\n",
			usage.Namespace,
			usage.PodName,
			usage.CPURequest,
			usage.CPULimit,
			usage.CPUUsage,
			cpuRateStr,
			usage.MemRequest,
			usage.MemLimit,
			usage.MemUsage,
			memRateStr,
			usage.PodStatus,
		)
	}
}

// 辅助函数：int64指针转换
func int64Ptr(i int64) *int64 {
	return &i
}
