package Model

type Image struct {
	Cluster      string `json:"cluster"`
	ImageName    string `json:"image_name"`
	PodName      string `json:"pod_name"`
	Namespace    string `json:"namespace"`
	WorkloadName string `json:"workload_name"`
	WorkloadType string `json:"workload_type"`
}
