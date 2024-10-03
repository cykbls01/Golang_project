package acs

import (
	"basic/util"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"log"
)

var Client *sdk.Client

func init() {
	var err error
	Client, err = sdk.NewClientWithOptions(util.Config["acr-regionid"], sdk.NewConfig(), credentials.NewAccessKeyCredential(util.Config["acs-ak"], util.Config["acs-sk"]))
	log.Println("init acs success")
	if err != nil {
		panic(err)
	}
}

func Request() *requests.CommonRequest {
	request := requests.NewCommonRequest()
	request.SetScheme(requests.HTTP)
	request.Product = "cr-ee"
	request.Version = "2018-12-01"
	request.Domain = util.Config["acr-endpoint"]
	request.Method = "POST"

	request.Headers["x-acs-caller-sdk-source"] = "cyk"
	request.Headers["x-acs-organizationid"] = util.Config["acr-organizationid"]
	request.Headers["x-acs-resourcegroupid"] = util.Config["acr-resourcegroupid"]
	request.Headers["x-acs-instanceid"] = util.Config["acr-instanceid"]
	request.Headers["x-acs-regionid"] = util.Config["acr-regionid"]
	return request
}
