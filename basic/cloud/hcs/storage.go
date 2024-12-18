package hcs

import (
	"basic/util"
)

type FileSystem struct {
	Size int    `json:"size"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

func ListFsDetail(projectId string) []FileSystem {
	body := Call(projectId, util.Config.MP["sfs-endpoint"]+"/v2/"+projectId+"/shares/detail", "GET", []byte{})
	type Resp struct {
		Shares []FileSystem `json:"shares"`
	}
	var rp Resp
	rp, _ = util.ParseJSON[Resp](body)
	return rp.Shares
}

func ListAllFsDetail(projectId []string) []FileSystem {
	res := make([]FileSystem, 0)
	for _, pid := range projectId {
		res = append(res, ListFsDetail(pid)...)
	}
	return res
}

func ListBlockDetail() {

}

func ListObsDetail() {

}
