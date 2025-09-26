package util

import (
	"flag"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var Config struct {
	MP       map[string]string
	Regions  map[string]Region
	Method   string
	DataPath string
}

type Repository struct {
	Namespace string `json:"RepoNamespaceName"`
	Name      string `json:"RepoName"`
	Id        string `json:"RepoId"`
}

type Tag struct {
	Namespace string `json:"Namespace"`
	Repo      string `json:"Repo"`
	Digest    string `json:"Digest"`
	Tag       string `json:"Tag"`
	Updated   string `json:"Updated"`
}

type Region struct {
	OrganizationID  string `yaml:"organizationid"`
	ResourceGroupID string `yaml:"resourcegroupid"`
	InstanceID      string `yaml:"instanceid"`
	RegionID        string `yaml:"regionid"`
	Filter          string `yaml:"filter"`
	Endpoint        string `yaml:"endpoint"`
	Registry        string `yaml:"registry"`
	Ak              string `yaml:"ak"`
	Sk              string `yaml:"sk"`
	Platform        string `yaml:"platform"`
}

func parseConfig(configPath string) map[string]string {
	dataBytes, err := os.ReadFile(configPath)
	mp := make(map[string]string)
	err = yaml.Unmarshal(dataBytes, mp)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(mp)
	return mp
}

func ParseRegion(filePath string) map[string]Region {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	var result map[string]Region
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return result
}

func Init() {
	var configPath string
	var regionPath string
	flag.StringVar(&configPath, "config", "config.yaml", "配置文件的路径")
	flag.StringVar(&regionPath, "region", "region.yaml", "区域文件的路径")
	flag.StringVar(&Config.Method, "main", "image_sync", "执行方法")
	flag.StringVar(&Config.DataPath, "data", "data.json", "数据文件")
	flag.Parse()
	Config.MP = parseConfig(configPath)
	Config.Regions = ParseRegion(regionPath)
}
