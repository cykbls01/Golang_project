package util

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var Config map[string]string

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

func Load(configPath string) map[string]string {
	dataBytes, err := os.ReadFile(configPath)
	mp := make(map[string]string)
	err = yaml.Unmarshal(dataBytes, mp)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(mp)
	return mp
}

func init() {

}
