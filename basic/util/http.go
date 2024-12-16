package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func Call(data, endpoint string) *http.Response {
	body, err := ioutil.ReadFile(data)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// 解析JSON数据到结构体（假设结构体为RequestData）
	var requestData map[string]interface{}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	// 将结构体转换回JSON，作为请求体
	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		log.Fatalf("Error marshaling request data: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// 设置请求头（如果需要）
	req.Header.Set("Content-Type", "application/json")

	// 忽略SSL证书验证（不推荐在生产环境中使用，仅用于示例）
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 禁用SSL证书验证
		},
	}

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()
	return resp
	// 读取响应内容
	//responseBody, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(resp.Header)
	//if err != nil {
	//	log.Fatalf("Error reading response body: %v", err)
	//}
	//
	//// 解析响应数据到结构体（假设结构体为ResponseData）
	//var responseData map[string]interface{}
	//err = json.Unmarshal(responseBody, &responseData)
	//if err != nil {
	//	log.Fatalf("Error unmarshaling response data: %v", err)
	//}
	//
	//// 输出或处理响应数据
	//log.Printf("Response data: %+v", responseData)
}
