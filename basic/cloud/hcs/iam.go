package hcs

import (
	"basic/util"
	"encoding/json"
	"fmt"
)

type Project struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func ListProject() []Project {
	token := GetGlobalToken()
	_, body := util.Call(util.Config.MP["hcs-project-url"], "GET", []byte{}, map[string]string{"X-Auth-Token": token})
	type Resp struct {
		Projects []Project `json:"projects"`
	}
	var rp Resp
	rp, _ = util.ParseJSON[Resp](body)
	return rp.Projects
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
