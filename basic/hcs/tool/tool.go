package tool

import (
	"basic/util"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/model"
	"time"
)

func FilterRepoByTime(list []model.ShowReposResp, offset int) []model.ShowReposResp {
	repos := make([]model.ShowReposResp, 0)
	for _, v := range list {
		if util.JudgeTimeWithUTC(v.UpdatedAt, time.Duration(offset)*time.Hour) {
			repos = append(repos, v)
		}
	}
	return repos
}

func FilterTagByTime(list []util.Tag, offset int) []util.Tag {
	repos := make([]util.Tag, 0)
	for _, v := range list {
		if util.JudgeTimeWithUTC(v.Updated, time.Duration(offset)*time.Hour) {
			repos = append(repos, v)
		}
	}
	return repos
}
