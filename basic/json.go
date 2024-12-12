package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"io/ioutil"
)

type Content struct {
	Content string
}

func jsonFileToStruct(filename string, v interface{}) error {
	// 读取文件内容
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// 解析JSON数据到指定的结构体实例
	if err := json.Unmarshal(jsonData, v); err != nil {
		return err
	}

	return nil
}

func main() {
	// 假设你的JSON文件名为 person.json，并且位于同一目录下
	filename := "person.json"

	// 创建一个Person结构体实例用于填充数据
	request := requests.NewCommonRequest()
	// 调用函数，从文件中读取JSON并转换为Person结构体实例
	if err := jsonFileToStruct(filename, &request); err != nil {
		fmt.Println("Error reading or parsing JSON:", err)
	}
	var content = Content{}
	jsonData, _ := ioutil.ReadFile(filename)
	json.Unmarshal(jsonData, &content)
	request.Content = []byte(content.Content)
	fmt.Println(request)
}
