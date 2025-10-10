package main

import (
	processor1 "basic/processor"
	"basic/util"
	_ "basic/util/k8s"
	"encoding/json"
	"github.com/kr/pretty"
	_ "github.com/kr/pretty"
	_ "k8s.io/api/core/v1"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
	var processor processor1.Processor
	switch util.Config.Method {
	case "image_check":
		{
			files, _ := ListAllFiles(util.Config.MP["path"])
			processor = &processor1.ImageCheck{Files: files, Path: "image_check"}
		}
	case "pod_check":
		{
			files, _ := ListAllFiles(util.Config.MP["path"])
			processor = &processor1.PodCheck{Files: files, Path: "pod_check", Result: make([]processor1.Data, 0)}
		}
	case "image_sync":
		{
			os.Create("images.yaml")
			for _, v := range strings.Split(util.Config.MP["kv"], "|") {
				log.Println(v)
				source := strings.Split(v, ":")[0]
				target := strings.Split(v, ":")[1]
				processor = &processor1.ImageSync{Source: util.Config.Regions[source], Target: util.Config.Regions[target]}
				processor.Pre()
				processor.Process()
				processor.Post()
			}
			return
		}
	//case "image_transfer":
	//	{
	//		processor = &processor1.ImageTransfer{}
	//	}
	default:
		pretty.Println("未知方法")
		os.Exit(1)
	}
	processor.Pre()
	processor.Process()
	processor.Post()

}
