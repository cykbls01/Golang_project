package hcs

import (
	"basic/util"
	"fmt"
)

type FileSystem struct {
	Size int    `json:"size"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

func ListFsDetail(projectId string) {
	body := Call(projectId, util.Config.MP["sfs-endpoint"]+"/v2/"+projectId+"/shares/detail", "GET", []byte{})
	type Resp struct {
		Shares []FileSystem `json:"shares"`
	}
	var rp Resp
	rp, _ = util.ParseJSON[Resp](body)
	fmt.Println(rp)
}

func ListBlockDetail() {

}

func ListObsDetail() {

}
