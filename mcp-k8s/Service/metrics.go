package Service

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"math"
	"mcp-k8s/Model"
	"mcp-k8s/Repository/Pod"
	"sort"
)

func GetMetrics(filepath, namespace string) {
	pods, _ := Pod.Pods(filepath, namespace)
	podMetricsList := Pod.Metrics(filepath, namespace)

	metricsCache := make(map[string]metricsv1beta1.PodMetrics)
	for _, pm := range podMetricsList.Items {
		key := fmt.Sprintf("%s/%s", pm.Namespace, pm.Name)
		metricsCache[key] = pm
	}

	var usageList []Model.PodResourceUsage
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			continue
		}

		metricsKey := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
		podMetrics, hasMetrics := metricsCache[metricsKey]
		cpuRequest := int64(0)
		cpuLimit := int64(0)
		memRequest := int64(0)
		memLimit := int64(0)

		for _, container := range pod.Spec.Containers {
			if req, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
				cpuRequest += req.Value() * 1000
			}
			if lim, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
				cpuLimit += lim.Value() * 1000
			}

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
			cpuUsageRate = math.Round((float64(cpuUsage)/float64(cpuLimit))*100*100) / 100 // 保留2位小数
		}

		memUsageRate := 0.0
		if memLimit > 0 && hasMetrics {
			memUsageRate = math.Round((float64(memUsage)/float64(memLimit))*100*100) / 100
		}

		// 格式化输出（转换为人类可读单位）
		usageList = append(usageList, Model.PodResourceUsage{
			Namespace:    pod.Namespace,
			PodName:      pod.Name,
			CPURequest:   formatCPU(cpuRequest),
			CPULimit:     formatCPU(cpuLimit),
			CPUUsage:     formatCPU(cpuUsage),
			CPUUsageRate: cpuUsageRate,
			MemRequest:   formatMemory(memRequest),
			MemLimit:     formatMemory(memLimit),
			MemUsage:     formatMemory(memUsage),
			MemUsageRate: memUsageRate,
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
func printUsageTable(usageList []Model.PodResourceUsage) {
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
		)
	}
}
