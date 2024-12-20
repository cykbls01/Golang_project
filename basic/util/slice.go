package util

import (
	"encoding/json"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/model"
	"reflect"
	"strings"
)

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

func ParseJSON[T any](data []byte) (T, error) {
	var target T
	if err := json.Unmarshal(data, &target); err != nil {
		return target, err
	}
	return target, nil
}

//func ParseJSON[T any](data []byte) ([]T, error) {
//	// 创建一个切片来存储解码后的T类型的值
//	var result []T
//
//	// 使用json.Unmarshal将字节切片解码到目标类型的切片中
//	if err := json.Unmarshal(data, &result); err != nil {
//		return nil, err
//	}
//	return result, nil
//}

// FilterStructsByFieldName 函数接受一个切片和一个字段名，返回一个新的切片，其中只包含字段值包含指定子字符串的元素
func FilterStructsByFieldName(slice interface{}, fieldName, subStr string) interface{} {
	// 获取切片的反射类型和值
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		panic("input is not a slice")
	}

	// 创建一个新的切片来存储结果
	resultSlice := reflect.MakeSlice(sliceValue.Type(), 0, sliceValue.Len())

	// 遍历原始切片
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i)
		if elem.Kind() == reflect.Struct {
			// 获取结构体字段的反射值
			fieldValue := elem.FieldByName(fieldName)
			if fieldValue.IsValid() && fieldValue.Kind() == reflect.String && strings.Contains(fieldValue.String(), subStr) {
				// 如果字段值包含指定子字符串，则将该元素添加到结果切片中
				resultSlice = reflect.Append(resultSlice, elem)
			}
		}
	}

	return resultSlice.Interface()
}
