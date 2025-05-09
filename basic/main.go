package main

import (
	"basic/cloud/acs"
	"basic/cloud/acs/acr"
	"basic/cloud/hcs"
	swr2 "basic/cloud/hcs/swr"
	"basic/cloud/hcs/tool"
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

	util.Init()
	acs.Init()
	hcs.Init()
	switch util.Config.Method {
	case "sync_image":
		namespaceList := strings.Split(util.Config.MP["filter"], "|")
		tagList := make([]util.Tag, 0)
		for _, namespace := range namespaceList {
			swrRepos := swr2.ListRepoByNamespace(namespace)
			acrRepos := acr.ListRepoByNamespace(namespace)
			list := util.FindUniqueRepo(swrRepos, acrRepos)
			for _, repo := range list {
				log.Println("repo create: " + repo.Namespace + "-" + repo.Name)
				acr.CreateRepo(repo)
			}
			offset, _ := strconv.Atoi(util.Config.MP["time-compare"])
			for _, v := range tool.FilterRepoByTime(swrRepos, offset) {
				swrTags := tool.FilterTagByTime(swr2.ListTagByRepo(v), offset)
				tagList = append(tagList, swrTags...)
			}
		}
		tagList = append(tagList, util.Tag{Namespace: "acr-test", Tag: "latest", Repo: "nginx"})
		log.Println("sync number: " + strconv.Itoa(len(tagList)))
		util.Write(tagList, util.Config.MP)
	case "backup_sync":
		namespaceList := strings.Split(util.Config.MP["filter"], "|")
		repoList := make([]util.Repository, 0)
		for _, namespace := range namespaceList {
			swrRepos := swr2.ListRepoByNamespace(namespace)
			for _, repo := range swrRepos {
				repoList = append(repoList, util.Repository{Name: repo.Name, Namespace: repo.Namespace})
			}
		}
		util.WriteRepo(repoList, util.Config.MP)
	default:
		fmt.Printf("Error: unknown method '%s' (should not happen with default)\n", util.Config.Method)
		flag.Usage()
		os.Exit(1)
	}
	log.Println("end")
}
