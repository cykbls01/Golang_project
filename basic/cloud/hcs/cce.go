package hcs

import (
	"basic/util"
	"encoding/json"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cce/v3/model"
)

func ListCluster(projectId string) []model.Cluster {
	body := Call(projectId, util.Config.MP["cce-endpoint"]+"/api/v3/projects/"+projectId+"/clusters", "GET", []byte{})
	var rp model.ListClustersResponse
	rp, _ = util.ParseJSON[model.ListClustersResponse](body)
	return *rp.Items
}

func GetKubeConfig(projectId, clusterId string) model.CreateKubernetesClusterCertResponse {
	type JSONData struct {
		Duration int `json:"duration"`
	}
	var jsonData JSONData
	jsonData.Duration = 1
	data, _ := json.Marshal(jsonData)
	body := Call(projectId, util.Config.MP["cce-endpoint"]+"/api/v3/projects/"+projectId+"/clusters/"+clusterId+"/clustercert", "POST", data)
	var rp model.CreateKubernetesClusterCertResponse
	rp, _ = util.ParseJSON[model.CreateKubernetesClusterCertResponse](body)
	return rp
}
