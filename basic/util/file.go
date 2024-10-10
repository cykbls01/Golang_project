package util

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"log"
	"os"
)

func Write(list []Tag, mp map[string]string, ctx context.Context) {
	tracer := otel.Tracer("image-syncer")
	_, span := tracer.Start(ctx, "file")
	defer span.End()

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
