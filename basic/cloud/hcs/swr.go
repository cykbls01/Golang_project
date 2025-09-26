package hcs

import (
	"basic/util"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	swr "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/model"
	"log"
	"strconv"
)

type SWR struct {
	Client *swr.SwrClient
	Region util.Region
}

func (hcs *SWR) Init() error {
	auth := basic.NewCredentialsBuilder().
		WithAk(hcs.Region.Ak).
		WithSk(hcs.Region.Sk).
		Build()
	hcs.Client = swr.NewSwrClient(
		swr.SwrClientBuilder().
			WithEndpoints([]string{hcs.Region.Endpoint}).
			WithCredential(auth).
			WithHttpConfig(config.DefaultHttpConfig().WithIgnoreSSLVerification(true)).
			Build())
	log.Println("init swr success")
	return nil
}

func (hcs *SWR) ListTagByRepo(repo model.ShowReposResp) []util.Tag {

	list := make([]util.Tag, 0)
	for i := 0; i < 1; i++ {
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

func (hcs *SWR) ListRepoByNamespace(namespace string) ([]model.ShowReposResp, error) {

	list := make([]model.ShowReposResp, 0)
	for i := 0; i < 1; i++ {
		limitRequest := "50"
		orderColumn := "updated_time"
		orderType := "desc"
		offsetRequest := strconv.Itoa(i * 30)

		request := &model.ListReposDetailsRequest{}
		request.Namespace = &namespace
		request.Limit = &limitRequest
		request.Offset = &offsetRequest
		request.OrderColumn = &orderColumn
		request.OrderType = &orderType
		response, err := hcs.Client.ListReposDetails(request)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		if len(*response.Body) == 0 {
			break
		}
		list = append(list, *response.Body...)
	}
	return list, nil
}
