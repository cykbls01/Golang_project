package acs

import (
	"basic/util"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/model"
	"log"
)

type ACR struct {
	Client *sdk.Client
	Region util.Region
}

func (acs *ACR) Init() error {
	var err error
	acs.Client, err = sdk.NewClientWithOptions(acs.Region.RegionID, sdk.NewConfig(), credentials.NewAccessKeyCredential(acs.Region.Ak, acs.Region.Sk))
	if err != nil {
		return err
	}
	log.Println("init acs success")
	return nil
}

func (acs *ACR) ListTagByRepo(repo util.Repository) []util.Tag {

	request := acs.AcrRequest()

	request.ApiName = "ListRepoTag"
	request.QueryParams["InstanceId"] = "cri-private"
	request.QueryParams["RepoId"] = repo.Id
	request.QueryParams["PageSize"] = "5"
	request.SetContentType(requests.Form)

	response, err := acs.Client.ProcessCommonRequest(request)
	if err != nil {
		log.Println(err.Error())
	}

	var data struct {
		Tags []util.Tag `json:"Images"`
	}
	fmt.Println(response.GetHttpContentString())
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

func (acs *ACR) CreateRepo(repo model.ShowReposResp) error {
	request := acs.AcrRequest()

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
		return err
	}
	log.Println(response)
	return nil
}

func (acs *ACR) ListRepoByNamespace(namespace string) ([]util.Repository, error) {
	// 创建API请求
	request := acs.AcrRequest()

	request.ApiName = "ListRepository"
	request.QueryParams["RepoNamespaceName"] = namespace
	request.QueryParams["InstanceId"] = "cri-private"
	request.SetContentType(requests.Form)

	response, err := acs.Client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}

	var data struct {
		Repositories []util.Repository `json:"Repositories"`
	}
	err = json.Unmarshal([]byte(response.GetHttpContentString()), &data)
	if err != nil {
		return nil, err
	}
	return data.Repositories, nil
}

func (acs *ACR) AcrRequest() *requests.CommonRequest {
	request := requests.NewCommonRequest()
	request.SetScheme(requests.HTTP)
	request.Product = "cr-ee"
	request.Version = "2018-12-01"
	request.Domain = acs.Region.Endpoint
	request.Method = "POST"

	request.Headers["x-acs-caller-sdk-source"] = "cyk"
	request.Headers["x-acs-organizationid"] = acs.Region.OrganizationID
	request.Headers["x-acs-resourcegroupid"] = acs.Region.ResourceGroupID
	request.Headers["x-acs-instanceid"] = acs.Region.InstanceID
	request.Headers["x-acs-regionid"] = acs.Region.RegionID
	return request
}
