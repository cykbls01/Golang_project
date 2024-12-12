package acs

import (
	"basic/util"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"log"
)

var Client *sdk.Client

func Init() {
	var err error
	Client, err = sdk.NewClientWithOptions(util.Config.MP["acr-regionid"], sdk.NewConfig(), credentials.NewAccessKeyCredential(util.Config.MP["acs-ak"], util.Config.MP["acs-sk"]))
	log.Println("init acs success")
	if err != nil {
		panic(err)
	}
}

func AcrRequest() *requests.CommonRequest {
	request := requests.NewCommonRequest()
	request.SetScheme(requests.HTTP)
	request.Product = "cr-ee"
	request.Version = "2018-12-01"
	request.Domain = util.Config.MP["acr-endpoint"]
	request.Method = "POST"

	request.Headers["x-acs-caller-sdk-source"] = "cyk"
	request.Headers["x-acs-organizationid"] = util.Config.MP["acr-organizationid"]
	request.Headers["x-acs-resourcegroupid"] = util.Config.MP["acr-resourcegroupid"]
	request.Headers["x-acs-instanceid"] = util.Config.MP["acr-instanceid"]
	request.Headers["x-acs-regionid"] = util.Config.MP["acr-regionid"]
	return request
}

func AckRequest() *requests.CommonRequest {
	request := requests.NewCommonRequest()
	request.SetScheme(requests.HTTP)
	request.Domain = util.Config.MP["ack-endpoint"]
	request.Method = "POST"
	request.Version = "2015-12-15"

	request.Headers["Content-Type"] = "application/json"
	request.Headers["x-acs-caller-sdk-source"] = "cyk"
	return request
}
