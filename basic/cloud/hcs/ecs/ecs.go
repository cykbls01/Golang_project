package ecs

import (
	"basic/util"
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

func ListServerID() []string {
	var ids []string
	for _, item := range ListServer("d98fcad2881a42b393894817cc46fddf") {
		ids = append(ids, item.Id)
	}
	return ids
}
