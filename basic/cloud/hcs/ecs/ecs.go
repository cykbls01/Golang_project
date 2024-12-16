package ecs

import (
	"basic/util"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
)

func ListServer(project_id string) []model.NovaServer {
	auth := basic.NewCredentialsBuilder().
		WithAk(util.Config.MP["hcs-ak"]).
		WithSk(util.Config.MP["hcs-sk"]).
		WithProjectId(project_id).
		Build()

	client := ecs.NewEcsClient(
		ecs.EcsClientBuilder().
			WithEndpoints([]string{util.Config.MP["ecs-endpoint"]}).
			WithCredential(auth).
			WithHttpConfig(config.DefaultHttpConfig().WithIgnoreSSLVerification(true)).
			Build())
	name := "cce"
	request := &model.NovaListServersDetailsRequest{Name: &name}
	response, _ := client.NovaListServersDetails(request)
	return *response.Servers
}

func ListServerID(projectIds []string) []string {
	var ids []string
	fmt.Println(projectIds)
	for _, pid := range projectIds {
		if pid == "" {
			continue
		}
		for _, item := range ListServer(pid) {
			ids = append(ids, item.Id)
		}
	}
	fmt.Println(ids)
	return ids
}
