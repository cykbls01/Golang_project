package mcp

import (
	"fmt"
	"go.uber.org/zap"
	"mcp-k8s/Model"
	"mcp-k8s/Repository/Cluster"
	"mcp-k8s/Util"
	"strings"
)

func GetEndpoints(filepath, namespace string) []Model.Service {
	res := make([]Model.Service, 0)
	cluster, err := Cluster.Build(filepath)
	if err != nil {
		Util.Logger.Error("构建集群失败", zap.Error(err))
		return []Model.Service{}
	}

	services, err := cluster.Services(namespace)
	if err != nil {
		Util.Logger.Error("获取服务列表失败", zap.Error(err))
		return []Model.Service{}
	}
	for _, service := range services.Items {

		query := make([]string, 0)
		for k, v := range service.Spec.Selector {
			query = append(query, fmt.Sprintf("%s=%s", k, v))
		}

		svc := Model.Service{
			Namespace:   service.Namespace,
			Port:        Util.PluckAndJoin(service.Spec.Ports, "Port", ","),
			ServiceID:   string(service.UID),
			ServiceName: service.Name,
			ServiceType: string(service.Spec.Type),
			NodePort:    Util.PluckAndJoin(service.Spec.Ports, "NodePort", ","),
			TargetPort:  Util.PluckAndJoinNested(service.Spec.Ports, "TargetPort.IntVal", ","),
			ClusterName: Cluster.GetName(filepath),
			CreateTime:  service.CreationTimestamp.Time,
		}

		pod, _ := cluster.PodByLabel(service.Namespace, strings.Join(query, ","))
		if len(pod.Items) != 0 && len(pod.Items[0].OwnerReferences) != 0 {
			svc.WorkloadName = pod.Items[0].OwnerReferences[0].Name
			svc.WorkloadType = pod.Items[0].OwnerReferences[0].Kind
		}

		switch service.Spec.Type {
		case "NodePort":
			{
				nodes := GetNodes(filepath)
				svc.IP = Util.PluckAndJoin(nodes, "NodeIP", ",")
				res = append(res, svc)
			}
		case "LoadBalancer":
			{
				svc.IP = Util.PluckAndJoin(service.Status.LoadBalancer.Ingress, "IP", ",")
				res = append(res, svc)
			}
		default:
			{
				break
			}
		}
	}
	return res
}

func McpGetEndpoints(dirPath string) []Model.Service {
	res := make([]Model.Service, 0)
	files, _ := Util.WalkAllFilesAbs(dirPath)
	for _, file := range files {
		Util.Logger.Info("process", zap.String("file", file))
		res = append(res, GetEndpoints(file, "")...)
	}
	return res
}
