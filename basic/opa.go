package main

import (
	"basic/processor"
	cks2 "basic/processor/cks"
	"basic/util"
	_ "basic/util/k8s"
	"encoding/json"
	"github.com/kr/pretty"
	_ "github.com/kr/pretty"
	_ "k8s.io/api/core/v1"
	"os"
	"path/filepath"
	"reflect"
)

func ListAllFiles(dirPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录，只添加文件
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func WriteSliceToFile(data interface{}, filename string) error {
	// 类型校验：必须为切片类型
	if reflect.ValueOf(data).Kind() != reflect.Slice {
		return os.ErrInvalid
	}

	// 序列化为带缩进的JSON格式
	formattedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// 原子化写入文件（覆盖模式+自动创建）
	return os.WriteFile(filename, formattedData, 0644)
}

type Count struct {
	Privileged int
	HostPaths  int
	Sysctls    int
}

var count Count

func main() {
	util.Init()
	files, _ := ListAllFiles(util.Config.MP["path"])
	pretty.Println(files)
	var processor processor.Processor
	switch util.Config.Method {
	case "image_check":
		{
			processor = &cks2.ImageCheck{Files: files, Path: "image_check"}
		}
	case "pod_check":
		{
			processor = &cks2.PodCheck{Files: files, Path: "pod_check", Result: make([]cks2.Data, 0)}
		}
	default:
		pretty.Println("未知方法")
		os.Exit(1)
	}
	processor.Pre()
	processor.Process()
	processor.Post()

}
