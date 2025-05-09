package util

import (
	"fmt"
	"log"
	"os"
)

func Write(list []Tag, mp map[string]string) {

	file, err := os.Create("images.yaml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// 遍历切片，将每个元素写入文件，每个元素后写入换行符
	for _, item := range list {
		link := "/" + item.Namespace + "/" + item.Repo + ":" + item.Tag
		line := mp["swr-image-endpoint"] + link + ": " + mp["acr-image-endpoint"] + link
		fmt.Fprintln(file, line)
	}
	log.Println("File written successfully.")
}

func WriteRepo(list []Repository, mp map[string]string) {

	file, err := os.Create("images.yaml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// 遍历切片，将每个元素写入文件，每个元素后写入换行符
	for _, item := range list {
		link := "/" + item.Namespace + "/" + item.Name
		line := mp["swr-image-endpoint"] + link + ": " + mp["acr-image-endpoint"] + link
		fmt.Fprintln(file, line)
	}
	log.Println("File written successfully.")
}
