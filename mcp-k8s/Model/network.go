package Model

type Endpoint struct {
	Namespace      string `json:"namespace"`
	IP             string `json:"ip"`
	Port           int32  `json:"port"`
	ServiceName    string `json:"service_name"`
	ServiceType    string `json:"service_type"`
	NodePort       int32  `json:"node_port"`
	TargetPort     int32  `json:"target_port"`
	DeploymentName string `json:"deployment_name"`
}
