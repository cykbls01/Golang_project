package processor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 定义JSON结构体（支持任意类型元素）
type Vulnerability struct {
	Id                string `json:"VulnerabilityID"`
	Cluster           string `json:"Cluster"`
	Namespace         string `json:"Namespace"`
	Image             string `json:"Image"`
	VulnerabilityType string `json:"Type"`
	PkgID             string `json:"PkgID"`
	PkgName           string `json:"PkgName"`
	InstalledVersion  string `json:"InstalledVersion"`
	FixedVersion      string `json:"FixedVersion"`
	Title             string `json:"Title"`
	Description       string `json:"Description"`
}

type Result struct {
	Vulnerability     []Vulnerability `json:"Vulnerabilities"`
	VulnerabilityType string          `json:"Type"`
}

type JSONData struct {
	Result       []Result `json:"Results"`
	ArtifactName string   `json:"ArtifactName"`
}

func isExcluded(ns string) bool {
	excludeList := []string{"kube-system", "arms-prom", "falco", "kube-node-lease", "kube-public", "istio-system"}
	for _, excluded := range excludeList {
		if excluded == ns {
			return true
		}
	}
	return false
}

// 递归处理目录的主函数
func ProcessJSONFiles(root string) []Vulnerability {
	//var allResults []Result
	var allData []Vulnerability
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("访问路径失败 %q: %v", path, err)
		}

		// 仅处理.json文件
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			data, err := parseJSONFile(path)
			if err != nil {
				return fmt.Errorf("解析文件失败 %q: %v", path, err)
			}
			//allResults = append(allResults, data.Result...)
			for _, v := range data.Result {
				for _, u := range v.Vulnerability {
					u.VulnerabilityType = v.VulnerabilityType
					u.Image = data.ArtifactName
					u.Namespace = strings.Split(path, "/")[len(strings.Split(path, "/"))-2]
					u.Cluster = strings.Split(path, "/")[len(strings.Split(path, "/"))-3]
					allData = append(allData, u)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(excel.Output(allData, output))
	return allData
}

// 解析单个JSON文件
func parseJSONFile(path string) (JSONData, error) {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		//return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	var jsonData JSONData
	if err := json.Unmarshal(fileContent, &jsonData); err != nil {
		//return nil, fmt.Errorf("JSON解析错误: %v", err)
	}

	return jsonData, nil
}
