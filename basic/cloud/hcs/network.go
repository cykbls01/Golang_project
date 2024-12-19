package hcs

import "basic/util"

type LoadBalancer struct {
	Vip  string `json:"vip_address"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

type VPC struct {
	Vip  string `json:"vip_address"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Subnet struct {
	Cidr string `json:"cidr"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

func ListLoadBalancer(projectId string) []LoadBalancer {
	body := Call(projectId, util.Config.MP["vpc-endpoint"]+"/v3/"+projectId+"/elb/loadbalancers", "GET", []byte{})
	type Resp struct {
		LoadBalancers []LoadBalancer `json:"loadbalancers"`
	}
	var rp Resp
	rp, _ = util.ParseJSON[Resp](body)
	return rp.LoadBalancers
}

func ListSubnet(projectId string) []Subnet {
	body := Call(projectId, util.Config.MP["vpc-endpoint"]+"/v1/"+projectId+"/subnets", "GET", []byte{})
	type Resp struct {
		Subnets []Subnet `json:"subnets"`
	}
	var rp Resp
	rp, _ = util.ParseJSON[Resp](body)
	return rp.Subnets
}

func ListVPC(projectId string) []VPC {
	body := Call(projectId, util.Config.MP["vpc-endpoint"]+"/v1/"+projectId+"/vpcs", "GET", []byte{})
	type Resp struct {
		VPC []VPC `json:"vpcs"`
	}
	var rp Resp
	rp, _ = util.ParseJSON[Resp](body)
	return rp.VPC
}

func ListAllLbDetail(projectId []string) []LoadBalancer {
	res := make([]LoadBalancer, 0)
	for _, pid := range projectId {
		res = append(res, ListLoadBalancer(pid)...)
	}
	return res
}
