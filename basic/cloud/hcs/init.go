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

func GetGlobalToken() string {
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
				Domain struct {
					Name string `json:"name"`
				} `json:"domain"`
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
					"domain": {
						"name": "中国人寿"
					}
				}
			}
		}`
	header, _ := util.Call(util.Config.MP["hcs-token-url"], "POST", []byte(jsonStr), map[string]string{})
	return header["X-Subject-Token"][0]
}

func GetProjectToken(projectId string) string {
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
	header, _ := util.Call(util.Config.MP["hcs-token-url"], "POST", data, map[string]string{})
	return header["X-Subject-Token"][0]
}

func Call(projectId, endpoint, method string, data []byte) []byte {
	_, body := util.Call(endpoint, method, data, map[string]string{
		"X-Auth-Token": GetProjectToken(projectId),
	})
	return body
}
