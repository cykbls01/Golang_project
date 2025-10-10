package processor

import (
	"basic/cloud/acs"
	"basic/cloud/hcs"
	"basic/cloud/hcs/tool"
	"basic/util"
	"log"
	"strconv"
	"strings"
)

type ImageSync struct {
	hcsClinet hcs.SWR
	Links     []string `json:"links"`
	Source    util.Region
	Target    util.Region
}

func (is *ImageSync) Pre() {
	is.hcsClinet = hcs.SWR{Region: is.Source}
	is.hcsClinet.Init()
}

func (is *ImageSync) Process() {

	switch is.Target.Platform {
	case "hcs":
		{
			tags := make([]util.Tag, 0)
			for _, namespace := range strings.Split(is.Target.Filter, "|") {
				log.Println(namespace)
				swrRepos, err := is.hcsClinet.ListRepoByNamespace(namespace)
				if err != nil {
					log.Println(err)
					continue
				}
				offset, _ := strconv.Atoi(util.Config.MP["time-compare"])
				for _, v := range tool.FilterRepoByTime(swrRepos, offset) {
					swrTags := tool.FilterTagByTime(is.hcsClinet.ListTagByRepo(v), offset)
					for _, u := range swrTags {
						log.Println("tag: " + u.Namespace + "/" + u.Repo + ":" + u.Tag)
						tags = append(tags, u)
					}
				}
			}

			for _, v := range tags {
				first, second := TruncateLastHyphen(v.Tag)
				link1 := "/" + v.Namespace + "/" + v.Repo + ":" + v.Tag
				link2 := "/" + v.Namespace + "/" + v.Repo + ":" + first
				var line string
				if ((second == "r5" || second == "beijing") && is.Target.RegionID == "r5") || ((second == "r7" || second == "shanghai") && is.Target.RegionID == "r7") || (second == "r9" && is.Target.RegionID == "r9") || (second == "both" && (is.Target.RegionID == "r5" || is.Target.RegionID == "r7")) {
					line = is.Source.Registry + link1 + ": " + is.Target.Registry + link2
				} else {
					continue
				}
				is.Links = append(is.Links, line)
			}
		}
	case "acs":
		{
			acsClient := acs.ACR{Region: is.Target}
			err := acsClient.Init()
			if err != nil {
				log.Println(err)
				return
			}
			for _, namespace := range strings.Split(is.Target.Filter, "|") {
				swrRepos, err := is.hcsClinet.ListRepoByNamespace(namespace)
				if err != nil {
					log.Println(err)
					continue
				}
				acrRepos, err := acsClient.ListRepoByNamespace(namespace)
				if err != nil {
					log.Println(err)
					continue
				}

				list := util.FindUniqueRepo(swrRepos, acrRepos)
				for _, repo := range list {
					log.Println("repo create: " + repo.Namespace + "-" + repo.Name)
					err = acsClient.CreateRepo(repo)
					if err != nil {
						log.Println(err)
					}
				}

				offset, _ := strconv.Atoi(util.Config.MP["time-compare"])
				for _, v := range tool.FilterRepoByTime(swrRepos, offset) {
					swrTags := tool.FilterTagByTime(is.hcsClinet.ListTagByRepo(v), offset)
					for _, u := range swrTags {
						link := "/" + v.Namespace + "/" + u.Repo + ":" + u.Tag
						is.Links = append(is.Links, is.Source.Registry+link+": "+is.Target.Registry+link)
					}
				}
			}
		}
	default:
		{
			return
		}
	}
}

func (is *ImageSync) Post() {
	nginx := "/" + "acr-test" + "/" + "nginx" + ":" + "latest"
	is.Links = append(is.Links, is.Source.Registry+nginx+": "+is.Target.Registry+nginx)
	log.Println("sync number: " + strconv.Itoa(len(is.Links)))
	util.WriteLines(is.Links)
}

func TruncateLastHyphen(s string) (prefix, suffix string) {
	lastIndex := strings.LastIndex(s, "-")

	if lastIndex == -1 {
		return s, ""
	}
	prefix = s[:lastIndex]
	suffix = s[lastIndex+1:]
	return
}
