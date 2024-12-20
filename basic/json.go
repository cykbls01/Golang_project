package main

import (
	_ "basic/cloud/acs"
	"basic/cloud/hcs"
	"basic/util"
	_ "basic/util"
	"encoding/json"
	"fmt"
	_ "github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"io/ioutil"
	"os"
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
	//util.Init()
	//acs.Init()
	//
	//request := requests.NewCommonRequest()
	//if err := jsonFileToStruct(util.Config.DataPath, &request); err != nil {
	//	fmt.Println("Error reading or parsing JSON:", err)
	//}
	//var content = Content{}
	//jsonData, _ := ioutil.ReadFile(util.Config.DataPath)
	//json.Unmarshal(jsonData, &content)
	//request.Content = []byte(content.Content)
	//fmt.Println(request)
	//response, err := acs.Client.ProcessCommonRequest(request)
	//if err != nil {
	//	log.Println(err.Error())
	//}
	//log.Println(response)
	util.Init()
	//excel.Output(hcs.ListAllLoadBalancer(hcs.ListProject()), "lb.xlsx")
	//excel.Output(hcs.ListAllVPC(hcs.ListProject()), "vpc.xlsx")
	//excel.Output(hcs.ListAllSubnet(hcs.ListProject()), "subnet.xlsx")
	for _, p := range hcs.ListProject() {
		for _, c := range hcs.ListCluster(p.Id) {
			config := hcs.GetKubeConfig(p.Id, *c.Metadata.Uid)
			jsonData, err := json.Marshal(config)
			if err != nil {
				fmt.Println("JSON marshaling failed:", err)
				return
			}

			// 将JSON数据写入文件
			err = os.WriteFile(c.Metadata.Name+".json", jsonData, 0644)
			if err != nil {
				fmt.Println("Failed to write JSON data to file:", err)
				return
			}
			break
		}
	}
}
