package Model

type PodResourceUsage struct {
	Namespace    string  `json:"namespace"`
	PodName      string  `json:"podName"`
	CPURequest   string  `json:"cpu_request"`
	CPULimit     string  `json:"cpu_limit"`
	CPUUsage     string  `json:"cpu_usage"`
	CPUUsageRate float64 `json:"cpu_usage_rate"`
	MemRequest   string  `json:"mem_request"`
	MemLimit     string  `json:"mem_limit"`
	MemUsage     string  `json:"mem_usage"`
	MemUsageRate float64 `json:"mem_usage_rate"`
}
