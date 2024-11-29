package acr

import (
	"basic/cloud/acs"
	"basic/util"
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/model"
	"log"
)

func CreateRepo(repo model.ShowReposResp) {
	request := acs.Request()

	// 接口业务参数设置
	request.ApiName = "CreateRepository"
	request.QueryParams["RepoNamespaceName"] = repo.Namespace
	request.QueryParams["RepoType"] = "PUBLIC"
	request.QueryParams["InstanceId"] = "cri-private"
	request.QueryParams["RepoName"] = repo.Name
	request.QueryParams["Summary"] = repo.Name
	request.SetContentType(requests.Form)

	response, err := acs.Client.ProcessCommonRequest(request)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(response)
}

func ListRepoByNamespace(namespace string) []util.Repository {
	// 创建API请求
	request := acs.Request()

	request.ApiName = "ListRepository"
	request.QueryParams["RepoNamespaceName"] = namespace
	request.QueryParams["InstanceId"] = "cri-private"
	request.SetContentType(requests.Form)

	response, err := acs.Client.ProcessCommonRequest(request)
	if err != nil {
		log.Println(err.Error())
	}

	var data struct {
		Repositories []util.Repository `json:"Repositories"`
	}
	err = json.Unmarshal([]byte(response.GetHttpContentString()), &data)
	if err != nil {
		log.Println(err.Error())
	}
	return data.Repositories
}
