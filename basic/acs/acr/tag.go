package acr

import (
	"basic/acs"
	"basic/util"
	"encoding/json"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"log"
)

func ListTagByRepo(repo util.Repository) []util.Tag {

	request := acs.Request()

	request.ApiName = "ListRepoTag"
	request.QueryParams["InstanceId"] = "cri-private"
	request.QueryParams["RepoId"] = repo.Id
	request.QueryParams["PageSize"] = "150"
	request.SetContentType(requests.Form)

	response, err := acs.Client.ProcessCommonRequest(request)
	if err != nil {
		log.Println(err.Error())
	}

	var data struct {
		Tags []util.Tag `json:"Images"`
	}
	err = json.Unmarshal([]byte(response.GetHttpContentString()), &data)
	if err != nil {
		log.Println(err.Error())
	}
	for k, _ := range data.Tags {
		data.Tags[k].Namespace = repo.Namespace
		data.Tags[k].Repo = repo.Name
	}
	return data.Tags
}
