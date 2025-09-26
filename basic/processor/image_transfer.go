package processor

//
//import (
//	"basic/cloud/hcs"
//	"basic/cloud/hcs/tool"
//	"basic/util"
//	"log"
//	"strconv"
//	"strings"
//)
//
//type ImageTransfer struct {
//	Tags       []util.Tag `json:"tags"`
//	Namespaces []string   `json:"namespaces"`
//	hcsClinet  hcs.SWR
//	Source     util.Region
//	Target     util.Region
//}
//
//func (it *ImageTransfer) Pre() {
//	it.hcsClinet = hcs.SWR{Region: it.Source}
//	it.hcsClinet.Init()
//	it.Namespaces = strings.Split(util.Config.MP["filter"], "|")
//}
//
//func (it *ImageTransfer) Process() {
//	for _, namespace := range it.Namespaces {
//		swrRepos, err := it.hcsClinet.ListRepoByNamespace(namespace)
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//		offset, _ := strconv.Atoi(util.Config.MP["time-compare"])
//		for _, v := range tool.FilterRepoByTime(swrRepos, offset) {
//			swrTags := tool.FilterTagByTime(it.hcsClinet.ListTagByRepo(v), offset)
//			for _, u := range swrTags {
//				log.Println("tag: " + u.Namespace + "/" + u.Repo + ":" + u.Tag)
//				it.Tags = append(it.Tags, u)
//			}
//		}
//	}
//}
//
//func (it *ImageTransfer) Post() {
//	list := make([]string, 0)
//	for _, v := range it.Tags {
//		first, second := TruncateLastHyphen(v.Tag)
//		link1 := "/" + v.Namespace + "/" + v.Repo + ":" + v.Tag
//		link2 := "/" + v.Namespace + "/" + v.Repo + ":" + first
//		var line string
//		if second == "r5" || second == "beijing" {
//			line = util.Config.MP["swr-r4-image-endpoint"] + link1 + ": " + util.Config.MP["swr-r5-image-endpoint"] + link2
//		} else if second == "r7" || second == "shanghai" {
//			line = util.Config.MP["swr-r4-image-endpoint"] + link1 + ": " + util.Config.MP["swr-r7-image-endpoint"] + link2
//		} else if second == "r9" {
//			line = util.Config.MP["swr-r4-image-endpoint"] + link1 + ": " + util.Config.MP["swr-r9-image-endpoint"] + link2
//		} else if second == "both" {
//			list = append(list, util.Config.MP["swr-r4-image-endpoint"]+link1+": "+util.Config.MP["swr-r5-image-endpoint"]+link2)
//			list = append(list, util.Config.MP["swr-r4-image-endpoint"]+link1+": "+util.Config.MP["swr-r7-image-endpoint"]+link2)
//			continue
//		} else {
//			continue
//		}
//		list = append(list, line)
//	}
//	nginx := "/" + "acr-test" + "/" + "nginx" + ":" + "latest"
//	list = append(list, util.Config.MP["swr-r4-image-endpoint"]+nginx+": "+util.Config.MP["swr-r4-image-endpoint"]+nginx)
//	log.Println("sync number: " + strconv.Itoa(len(list)))
//	util.WriteLines(list)
//}
