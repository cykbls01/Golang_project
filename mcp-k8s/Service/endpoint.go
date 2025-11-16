package Service

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	"log/slog"
	"mcp-k8s/Model"
	"mcp-k8s/Repository/Cluster"
	"mcp-k8s/Repository/Network"
	"mcp-k8s/Repository/Pod"
	"mcp-k8s/Util"
	"strings"
)

func GetEndpoints(filepath, namespace string) []Model.Endpoint {
	res := make([]Model.Endpoint, 0)

	services, err := Network.Services(filepath, namespace)
	if err != nil {
		slog.Error("获取服务列表失败", "error", err)
		return []Model.Endpoint{}
	}
	for _, service := range services.Items {

		query := make([]string, 0)
		for k, v := range service.Spec.Selector {
			query = append(query, fmt.Sprintf("%s=%s", k, v))
		}
		deployment, _ := Pod.DeploymentByLabel(filepath, service.Namespace, strings.Join(query, ","))
		deploymentName := ""
		if len(deployment.Items) != 0 {
			deploymentName = deployment.Items[0].Name
		}

		switch service.Spec.Type {
		case "NodePort":
			{
				res = append(res, NodePortService(filepath, service, deploymentName)...)
			}
		case "LoadBalancer":
			{
				res = append(res, LoadBalancerService(filepath, service, deploymentName)...)
			}
		default:
			{

			}
		}
	}
	return res
}

func NodePortService(filepath string, service v1.Service, deploymentName string) []Model.Endpoint {
	res := make([]Model.Endpoint, 0)
	for _, port := range service.Spec.Ports {
		nodes, _ := Cluster.Nodes(filepath)
		for _, v := range nodes.Items {
			res = append(res, Model.Endpoint{
				Namespace:      service.Namespace,
				IP:             v.Status.Addresses[0].Address,
				Port:           port.Port,
				ServiceName:    service.Name,
				ServiceType:    string(service.Spec.Type),
				NodePort:       port.NodePort,
				TargetPort:     port.TargetPort.IntVal,
				DeploymentName: deploymentName,
			})
		}
	}
	return res
}

func LoadBalancerService(filepath string, service v1.Service, deploymentName string) []Model.Endpoint {
	res := make([]Model.Endpoint, 0)
	for _, port := range service.Spec.Ports {
		res = append(res, Model.Endpoint{
			Namespace:      service.Namespace,
			IP:             service.Status.LoadBalancer.Ingress[0].IP,
			Port:           port.Port,
			ServiceName:    service.Name,
			ServiceType:    string(service.Spec.Type),
			NodePort:       port.NodePort,
			TargetPort:     port.TargetPort.IntVal,
			DeploymentName: deploymentName,
		})
	}
	return res
}

func McpGetEndpoints(dirPath string) []Model.Endpoint {
	res := make([]Model.Endpoint, 0)
	files, _ := Util.WalkAllFilesAbs(dirPath)
	for _, file := range files {
		res = append(res, GetEndpoints(file, "")...)
	}
	return res
}
