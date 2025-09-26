package hcs

import (
	"basic/util"
)

//func Init() *swr.SwrClient {
//	auth := basic.NewCredentialsBuilder().
//		WithAk(util.Config.MP["hcs-ak"]).
//		WithSk(util.Config.MP["hcs-sk"]).
//		Build()
//	Client := swr.NewSwrClient(
//		swr.SwrClientBuilder().
//			WithEndpoints([]string{util.Config.MP["swr-api-endpoint"]}).
//			WithCredential(auth).
//			WithHttpConfig(config.DefaultHttpConfig().WithIgnoreSSLVerification(true)).
//			Build())
//	log.Println("init swr success")
//	return Client
//}

func Call(projectId, endpoint, method string, data []byte) []byte {
	_, body := util.Call(endpoint, method, data, map[string]string{
		"X-Auth-Token": GetProjectToken(projectId),
		"Content-Type": "application/json",
	})
	return body
}
