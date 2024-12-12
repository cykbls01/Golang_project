package ack

import (
	"basic/cloud/acs"
	"log"
)

func DeployPolicy(id, policy string) {
	request := acs.AckRequest()
	request.PathPattern = "/clusters/" + id + "/policies/" + policy
	body := `{"action":"deny","parameters":{"allowedHostPaths":[{"readOnly":true,"pathPrefix":"/"}]}}`
	request.Content = []byte(body)

	response, err := acs.Client.ProcessCommonRequest(request)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(response)
}
