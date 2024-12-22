package k8s

import (
	"basic/cloud/hcs"
	"encoding/json"
	"fmt"
	"os"
)

func GenerateConfigs() {
	for _, p := range hcs.ListProject() {
		for _, c := range hcs.ListCluster(p.Id) {
			config := hcs.GetKubeConfig(p.Id, *c.Metadata.Uid)
			jsonData, err := json.Marshal(config)
			if err != nil {
				fmt.Println("JSON marshaling failed:", err)
				return
			}

			// 将JSON数据写入文件
			err = os.WriteFile(c.Metadata.Name+".json", jsonData, 0644)
			if err != nil {
				fmt.Println("Failed to write JSON data to file:", err)
				return
			}
			break
		}
	}
}
