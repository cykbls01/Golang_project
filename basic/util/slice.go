package util

import "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/model"

func Filter[T any](slice []T, predicate func(T) bool) []T {
	var filtered []T
	for _, item := range slice {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func FindUniqueRepo(list1 []model.ShowReposResp, list2 []Repository) []model.ShowReposResp {
	// 创建一个map来存储list2中所有元素的name
	list2Names := make(map[string]bool)
	for _, item := range list2 {
		list2Names[item.Name+"-"+item.Namespace] = true
	}

	// 遍历list1，查找不在list2中的元素
	var uniqueList []model.ShowReposResp
	for _, item := range list1 {
		if _, exists := list2Names[item.Name+"-"+item.Namespace]; !exists {
			uniqueList = append(uniqueList, item)
		}
	}

	return uniqueList
}

func FindUniqueTag(list1, list2 []Tag) []Tag {
	// 创建一个map来存储list2中所有元素的name
	list2Names := make(map[string]bool)
	for _, item := range list2 {
		list2Names[item.Namespace+"-"+item.Repo+"-"+item.Tag] = true
	}

	// 遍历list1，查找不在list2中的元素
	var uniqueList []Tag
	for _, item := range list1 {
		if _, exists := list2Names[item.Namespace+"-"+item.Repo+"-"+item.Tag]; !exists {
			uniqueList = append(uniqueList, item)
		}
	}

	return uniqueList
}
