package swr

import (
	"basic/hcs"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/model"
	"log"
	"strconv"
)

func ListRepoByNamespace(namespace string) []model.ShowReposResp {

	list := make([]model.ShowReposResp, 0)
	for i := 0; i < 100; i++ {
		request := &model.ListReposDetailsRequest{}
		request.Namespace = &namespace
		limitRequest := "30"
		request.Limit = &limitRequest
		offsetRequest := strconv.Itoa(i * 30)
		request.Offset = &offsetRequest
		response, err := hcs.Client.ListReposDetails(request)
		if err != nil {
			log.Println(err.Error())
			break
		}
		if len(*response.Body) == 0 {
			break
		}
		list = append(list, *response.Body...)
	}
	return list
}
