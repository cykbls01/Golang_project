package Entity

type Endpoint struct {
	Namespace   string
	IP          string
	Port        int32
	ServiceName string
	ServiceType string
	NodePort    int32
	TargetPort  int32
	PodName     string
}
