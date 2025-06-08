package cks

import (
	"basic/util"
	"basic/util/k8s"
	"bytes"
	"fmt"
	"github.com/kr/pretty"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ImageCheck struct {
	Files []string
	Path  string
	Error []string
}

func (ic *ImageCheck) Pre() {
	os.Mkdir(ic.Path, 0755)
}

func (ic *ImageCheck) Process() {
	for i, file := range ic.Files {
		ic.CheckImage(file, ic.Path)
		fmt.Printf("\rProcessing: %d/%d ", i+1, len(ic.Files))
	}
}

func (ic *ImageCheck) Post() {
	pretty.Println(ic.Error)
	ProcessJSONFiles(ic.Path, ic.Path+"/result.xlsx")
}

func ExecuteComplexCmd(cmdStr string) (string, string, error) {
	// 获取当前工作目录[4](@ref)
	//fmt.Println(cmdStr)
	currentDir, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("获取目录失败: %v", err)
	}

	// 创建带bash解析的命令对象[1,6](@ref)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	cmd.Dir = currentDir // 显式设置工作目录

	// 分离标准输出和错误缓冲区[3](@ref)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// 执行并等待命令完成[6](@ref)
	err = cmd.Run()
	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()

	return stdout, stderr, err
}

func (ic *ImageCheck) CheckImage(file, base string) {
	for _, ns := range k8s.ListAllNamespaces(file) {
		if isExcluded(ns.Namespace) {
			continue
		}
		os.Mkdir(base+"/"+filepath.Base(file), 0755)
		os.Mkdir(base+"/"+filepath.Base(file)+"/"+ns.Namespace, 0755)
		nsbase := base + "/" + filepath.Base(file) + "/" + ns.Namespace
		for _, image := range k8s.ListAllImages(k8s.ListAllPods(ns)) {
			//cmd1 := "trivy image --skip-db-update --skip-java-db-update  " + image + " --severity=CRITICAL  --report=summary --insecure --format template --template \"@/data/.cache/trivy/html.tpl\" --output " + nsbase + "/" + filepath.Base(image) + ".html --timeout 3m -q --cache-dir /data/.cache/trivy"
			cmd2 := "trivy image --skip-db-update --skip-java-db-update  " + image + " --severity=CRITICAL  --report=summary --insecure --format json --output " + nsbase + "/" + filepath.Base(image) + ".json --timeout 3m -q --cache-dir /data/.cache/trivy"
			switch {
			case strings.HasPrefix(image, "cr-ee"):
				//ExecuteComplexCmd(cmd1)
				ExecuteComplexCmd(cmd2)
			case strings.HasPrefix(image, "swr.bj-kjy-50"):
				//ExecuteComplexCmd(cmd1 + " --username " + util.Config.MP["swr5-username"] + " --password " + util.Config.MP["swr-password"])
				ExecuteComplexCmd(cmd2 + " --username " + util.Config.MP["swr5-username"] + " --password " + util.Config.MP["swr-password"])
			case strings.HasPrefix(image, "swr.bj-kjy-90"):
				//ExecuteComplexCmd(cmd1 + " --username " + util.Config.MP["swr9-username"] + " --password " + util.Config.MP["swr-password"])
				ExecuteComplexCmd(cmd2 + " --username " + util.Config.MP["swr9-username"] + " --password " + util.Config.MP["swr-password"])
			case strings.HasPrefix(image, "swr.sh-ky-70"):
				//ExecuteComplexCmd(cmd1 + " --username " + util.Config.MP["swr7-username"] + " --password " + util.Config.MP["swr-password"])
				ExecuteComplexCmd(cmd2 + " --username " + util.Config.MP["swr7-username"] + " --password " + util.Config.MP["swr-password"])
			default:
				ic.Error = append(ic.Error, image)
			}

		}
	}

}
