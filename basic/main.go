package main

import (
	"basic/acs/acr"
	"basic/hcs/swr"
	"basic/hcs/tool"
	"basic/util"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.Println("begin")
	var configPath string
	var method string
	flag.StringVar(&configPath, "config", "config.yaml", "配置文件的路径")
	flag.StringVar(&method, "method", "main", "The method to execute")
	flag.Parse()
	util.Config = util.Load(configPath)

	switch method {
	case "main":
		namespaceList := strings.Split(util.Config["filter"], "|")
		tagList := make([]util.Tag, 0)
		for _, namespace := range namespaceList {
			swrRepos := swr.ListRepoByNamespace(namespace)
			acrRepos := acr.ListRepoByNamespace(namespace)
			list := util.FindUniqueRepo(swrRepos, acrRepos)
			for _, repo := range list {
				log.Println("repo create: " + repo.Namespace + "-" + repo.Name)
				acr.CreateRepo(repo)
			}
			offset, _ := strconv.Atoi(util.Config["time-compare"])
			for _, v := range tool.FilterRepoByTime(swrRepos, offset) {
				swrTags := tool.FilterTagByTime(swr.ListTagByRepo(v), offset)
				tagList = append(tagList, swrTags...)
			}
		}
		util.Write(tagList, util.Config)
	default:
		fmt.Printf("Error: unknown method '%s' (should not happen with default)\n", method)
		flag.Usage()
		os.Exit(1)
	}
	log.Println("end")
}
