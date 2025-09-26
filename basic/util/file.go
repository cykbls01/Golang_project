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

	for _, item := range list {
		link := "/" + item.Namespace + "/" + item.Repo + ":" + item.Tag
		line := mp["source-endpoint"] + link + ": " + mp["target-endpoint"] + link
		fmt.Fprintln(file, line)
	}
	log.Println("File written successfully.")
}

func WriteLines(list []string) {
	file, err := os.OpenFile("images.yaml", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	for _, item := range list {
		_, err = fmt.Fprintln(file, item)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
	}
	log.Println("File written successfully.")
}

func WriteRepo(list []Repository, mp map[string]string) {

	file, err := os.Create("images.yaml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	for _, item := range list {
		link := "/" + item.Namespace + "/" + item.Name
		line := mp["swr-image-endpoint"] + link + ": " + mp["acr-image-endpoint"] + link
		fmt.Fprintln(file, line)
	}
	log.Println("File written successfully.")
}
