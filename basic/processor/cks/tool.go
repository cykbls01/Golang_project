package cks

import (
	"basic/excel"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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
func ProcessJSONFiles(root, output string) ([]Result, error) {
	var allResults []Result
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
			allResults = append(allResults, data.Result...)
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
	fmt.Println(excel.Output(allData, output))
	return allResults, err
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

func main() {
	// 获取当前工作目录
	currentDir, _ := os.Getwd()
	targetDir := filepath.Join(currentDir, os.Args[1])

	// 执行解析
	results, err := ProcessJSONFiles(targetDir, targetDir+"/result.xlsx")
	if err != nil {
		fmt.Printf("处理过程中发生错误: %v\n", err)
		return
	}

	vul := make([]Vulnerability, 0)
	for _, v := range results {
		for _, u := range v.Vulnerability {
			u.VulnerabilityType = v.VulnerabilityType
			vul = append(vul, u)
		}
	}

	typeMap := map[string]int{}
	idMap := map[string]int{}
	vulMap := make(map[string]Vulnerability)

	for _, v := range vul {
		typeMap[v.VulnerabilityType]++
		idMap[v.Id]++
		vulMap[v.Id] = v
	}

	fmt.Println("漏洞总数: " + strconv.Itoa(len(vul)))
	fmt.Println("漏洞总数分布: ")
	for k, v := range typeMap {
		fmt.Print(k)
		fmt.Print(" : ")
		fmt.Println(v)
	}

	typeMap = map[string]int{}
	for k, _ := range idMap {
		typeMap[vulMap[k].VulnerabilityType]++
	}
	fmt.Println("漏洞种类: " + strconv.Itoa(len(idMap)))
	fmt.Println("漏洞种类分布: ")
	for k, v := range typeMap {
		fmt.Print(k)
		fmt.Print(" : ")
		fmt.Println(v)
	}

}
