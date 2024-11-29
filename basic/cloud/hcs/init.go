package hcs

import (
	"basic/util"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	swr "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2"
	"log"
)

var Client *swr.SwrClient

func init() {
	auth := basic.NewCredentialsBuilder().
		WithAk(util.Config.MP["hcs-ak"]).
		WithSk(util.Config.MP["hcs-sk"]).
		Build()
	Client = swr.NewSwrClient(
		swr.SwrClientBuilder().
			WithEndpoints([]string{util.Config.MP["swr-endpoint"]}).
			WithCredential(auth).
			WithHttpConfig(config.DefaultHttpConfig().WithIgnoreSSLVerification(true)).
			Build())
	log.Println("init swr success")
}
