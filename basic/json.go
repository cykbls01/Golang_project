package main

import (
	"basic/cloud/acs"
	"basic/util"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"io/ioutil"
	"log"
)

type Content struct {
	Content string
}

func jsonFileToStruct(filename string, v interface{}) error {
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonData, v); err != nil {
		return err
	}
	return nil
}

func main() {
	util.Init()
	acs.Init()

	request := requests.NewCommonRequest()
	if err := jsonFileToStruct(util.Config.DataPath, &request); err != nil {
		fmt.Println("Error reading or parsing JSON:", err)
	}
	var content = Content{}
	jsonData, _ := ioutil.ReadFile(util.Config.DataPath)
	json.Unmarshal(jsonData, &content)
	request.Content = []byte(content.Content)
	fmt.Println(request)
	response, err := acs.Client.ProcessCommonRequest(request)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(response)
}
