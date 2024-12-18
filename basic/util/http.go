package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func CallFromFile(endpoint, method, data string, headers map[string]string) *http.Response {
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
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 禁用SSL证书验证
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()
	return resp
}

func Call(endpoint, method string, data []byte, headers map[string]string) (http.Header, []byte) {
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	fmt.Println(req)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 禁用SSL证书验证
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return resp.Header, body
}
