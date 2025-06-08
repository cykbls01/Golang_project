package main

import (
	"basic/cloud/acs"
	"basic/cloud/acs/acr"
	"basic/cloud/hcs"
	swr2 "basic/cloud/hcs/swr"
	"basic/util"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	log.Println("begin")

	util.Init()
	acs.Init()
	hcs.Init()
	switch util.Config.Method {
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
	case "test":
		fmt.Println(acr.ListTagByRepo(util.Repository{Name: "ezc", Namespace: "ezc", Id: "crr-oen7ow6qe4ay520f"}))
	default:
		fmt.Printf("Error: unknown method '%s' (should not happen with default)\n", util.Config.Method)
		flag.Usage()
		os.Exit(1)
	}
	log.Println("end")
}
