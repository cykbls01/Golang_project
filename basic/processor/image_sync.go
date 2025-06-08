package processor

import (
	"basic/cloud/acs/acr"
	swr2 "basic/cloud/hcs/swr"
	"basic/cloud/hcs/tool"
	"basic/util"
	"log"
	"strconv"
	"strings"
)

type ImageSync struct {
	Tags       []util.Tag `json:"tags"`
	Namespaces []string   `json:"namespaces"`
}

func (is *ImageSync) Pre() {
	is.Namespaces = strings.Split(util.Config.MP["filter"], "|")
}

func (is *ImageSync) Process() {
	for _, namespace := range is.Namespaces {
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
			is.Tags = append(is.Tags, swrTags...)
		}
	}
	is.Tags = append(is.Tags, util.Tag{Namespace: "acr-test", Tag: "latest", Repo: "nginx"})

}

func (is *ImageSync) Post() {
	log.Println("sync number: " + strconv.Itoa(len(is.Tags)))
	util.Write(is.Tags, util.Config.MP)
}
