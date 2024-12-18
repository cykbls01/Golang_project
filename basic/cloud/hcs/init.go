package hcs

import (
	"basic/util"
	"encoding/json"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	swr "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2"
	"log"
)

var Client *swr.SwrClient

func Init() {
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

func GetProjectToken(projectId string) string {
	return "MIIEZQYJKoZIhvcNAQcCoIIEVjCCBFICAQExDTALBglghkgBZQMEAgEwggLGBgkqhkiG9w0BBwGgggK3BIICs3sidG9rZW4iOnsiZXhwaXJlc19hdCI6IjIwMjQtMTItMTlUMDY6NDc6NDcuNzc3MDAwWiIsIm1ldGhvZHMiOlsicGFzc3dvcmQiXSwiY2F0YWxvZyI6W10sInJvbGVzIjpbeyJuYW1lIjoicmVhZG9ubHkiLCJpZCI6ImFhYmIzYjk4OGQxMDQ2ZjliOWE3NTAxYTZiYzgxZmFjIn0seyJuYW1lIjoiQ0NFX0Z1bGxBY2Nlc3MiLCJpZCI6ImM3OWY4ZTY3M2NmZTRhMmRiZjcxMGMzMTdmNjdlMjZjIn0seyJuYW1lIjoidGVfYWRtaW4iLCJpZCI6ImIzY2E1YWY2NDA2ZDQ1ZTJhZTA2OTkxZWJhMjIxNTI5In0seyJuYW1lIjoidmRjX3VzZXIiLCJpZCI6ImM4ZjY3YjY2YTNkZDQzZWViMTQ0OTkzYjFkZWJlYTI5In1dLCJwcm9qZWN0Ijp7ImRvbWFpbiI6eyJuYW1lIjoi5Lit5Zu95Lq65a+-IiwiaWQiOiI0ZjZiZjA1YmZhNzc0ZjhlOGVlOTg0Y2JhZWNiNzUwNCJ9LCJuYW1lIjoiYmota2p5LTUwX0dGLXNoYW5kb25nIiwiaWQiOiJkOGRjMTg5ZDM5ZmI0MDkwOTFhNjFlYjIyODQ3NjgwNSJ9LCJpc3N1ZWRfYXQiOiIyMDI0LTEyLTE4VDA2OjQ3OjQ3Ljc3NzAwMFoiLCJ1c2VyIjp7ImRvbWFpbiI6eyJuYW1lIjoi5Lit5Zu95Lq65a+-IiwiaWQiOiI0ZjZiZjA1YmZhNzc0ZjhlOGVlOTg0Y2JhZWNiNzUwNCJ9LCJuYW1lIjoiYWNzLXlpbmd5b25neXVuIiwiaWQiOiIyZmY2MDA1MzIzMmM0ZjdjOTRiMmY0YjI3OTM2NDYxYyJ9fX0xggFyMIIBbgIBATBJMD0xCzAJBgNVBAYTAkNOMQ8wDQYDVQQKEwZIdWF3ZWkxHTAbBgNVBAMTFEh1YXdlaSBJVCBQcm9kdWN0IENBAggVqZEhukGatzALBglghkgBZQMEAgEwDQYJKoZIhvcNAQEBBQAEggEASzojEZ1XQQqtRCAFKxqdeGHop2NCWeYInogMN7atAr0LNRCY1sriW5KU+Z54NjGfWIHdDerDvXOow+NXSvYqstRkdpp705lK567ii4LgJYnymnIzWdjYkdGzIOm92qtyd44eDiUyJK4y0R7vkTMt9E+PTaYHr+8zXGCoJ+O533r5wZgMVGKSnBTMEl1jj37pDYSfuDPlVozhoada+4KQbrAZ5-uAkQ8CFZbaMHC0z3iwo9P-WcDI9ygrXM52D8UzY3AhDRkz+loKIIj20ulkys25P1QZqaa2360vOMVP--KmOJveww2T3qsICpL7VV8WoUloug3ReEfg6K7ZX9PIgw=="
	type JSONData struct {
		Auth struct {
			Identity struct {
				Methods  []string `json:"methods"`
				Password struct {
					User struct {
						Name     string `json:"name"`
						Password string `json:"password"`
						Domain   struct {
							Name string `json:"name"`
						} `json:"domain"`
					} `json:"user"`
				} `json:"password"`
			} `json:"identity"`
			Scope struct {
				Project struct {
					ID string `json:"id"`
				} `json:"project"`
			} `json:"scope"`
		} `json:"auth"`
	}

	jsonStr := `{
			"auth": {
				"identity": {
					"methods": [
						"password"
					],
					"password": {
						"user": {
							"name": "acs-yingyongyun", 
							"password": "Kjy@2020",
							"domain": {
								"name": "中国人寿"
							}
						}
					}
				},
				"scope": {
					"project": {
						"id": "1234"
					}
				}
			}
		}`

	var jsonData JSONData
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
	}
	jsonData.Auth.Scope.Project.ID = projectId
	data, _ := json.Marshal(jsonData)
	fmt.Printf("%+v\n", jsonData)
	header, _ := util.Call(util.Config.MP["hcs-token-url"], "POST", data, map[string]string{})
	return header["X-Subject-Token"][0]
}

func Call(projectId, endpoint, method string, data []byte) []byte {
	_, body := util.Call(endpoint, method, data, map[string]string{
		"X-Auth-Token": GetProjectToken(projectId),
	})
	return body
}
