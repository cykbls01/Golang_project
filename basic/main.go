package main

import (
	"basic/acs/acr"
	"basic/hcs/swr"
	"basic/hcs/tool"
	"basic/util"
	"context"
	"flag"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func trace() context.Context {
	tracer := otel.Tracer("image-syncer")
	ctx, span := tracer.Start(context.Background(), "main")
	defer span.End()
	span.SetAttributes(attribute.String("event", "write called"))
	return ctx
}

func main() {
	log.Println("begin")
	switch util.Config.Method {
	case "main":
		ctx := trace()
		namespaceList := strings.Split(util.Config.MP["filter"], "|")
		tagList := make([]util.Tag, 0)
		for _, namespace := range namespaceList {
			swrRepos := swr.ListRepoByNamespace(namespace)
			acrRepos := acr.ListRepoByNamespace(namespace)
			list := util.FindUniqueRepo(swrRepos, acrRepos)
			for _, repo := range list {
				log.Println("repo create: " + repo.Namespace + "-" + repo.Name)
				acr.CreateRepo(repo)
			}
			offset, _ := strconv.Atoi(util.Config.MP["time-compare"])
			for _, v := range tool.FilterRepoByTime(swrRepos, offset) {
				swrTags := tool.FilterTagByTime(swr.ListTagByRepo(v), offset)
				tagList = append(tagList, swrTags...)
			}
		}
		util.Write(tagList, util.Config.MP, ctx)
	default:
		fmt.Printf("Error: unknown method '%s' (should not happen with default)\n", util.Config.Method)
		flag.Usage()
		os.Exit(1)
	}
	time.Sleep(5 * time.Second)
	log.Println("end")
}
