package main

//
//import (
//	"fmt"
//	"github.com/AliyunContainerService/image-syncer/pkg/client"
//	"github.com/AliyunContainerService/image-syncer/pkg/utils"
//)
//
//var (
//	logPath, configFile, authFile, imagesFile, successImagesFile, outputImagesFormat string
//
//	procNum, retries int
//
//	osFilterList, archFilterList []string
//
//	forceUpdate bool
//)
//
//func main() {
//
//	client, err := client.NewSyncClient("", "auth.yaml", "images.yaml", "", "", "yaml",
//		5, 2, utils.RemoveEmptyItems([]string{}), utils.RemoveEmptyItems([]string{}), false)
//	if err != nil {
//		fmt.Errorf("init sync client error: %v", err)
//		return
//	}
//
//	client.Run()
//
//}
