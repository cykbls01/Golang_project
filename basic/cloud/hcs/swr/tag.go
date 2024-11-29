package swr

import (
	"basic/cloud/hcs"
	"basic/util"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/model"
	"log"
	"strconv"
)

func ListTagByRepo(repo model.ShowReposResp) []util.Tag {

	list := make([]util.Tag, 0)
	for i := 0; i < 100; i++ {
		request := &model.ListRepositoryTagsRequest{}
		request.Namespace = repo.Namespace
		request.Repository = repo.Name
		limitRequest := "30"
		request.Limit = &limitRequest
		offsetRequest := strconv.Itoa(i * 30)
		request.Offset = &offsetRequest
		response, err := hcs.Client.ListRepositoryTags(request)
		if err != nil {
			log.Println(err.Error())
			break
		}
		if len(*response.Body) == 0 {
			break
		}
		for _, v := range *response.Body {
			list = append(list, util.Tag{Namespace: repo.Namespace, Repo: repo.Name, Tag: v.Tag, Updated: v.Updated})
		}
	}

	return list
}
