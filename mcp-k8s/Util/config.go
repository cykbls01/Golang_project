package Util

import (
	"fmt"
	"os"
	"path/filepath"
)

func WalkAllFilesAbs(dirPath string) ([]string, error) {
	// 1. 验证目录是否存在且合法
	info, err := os.Stat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("目录验证失败: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("路径 %s 不是有效目录", dirPath)
	}

	// 2. 存储绝对路径的切片
	var fileAbsPaths []string

	// 3. 递归遍历目录（核心逻辑）
	err = filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		// 处理遍历过程中的错误（如权限不足、符号链接失效）
		if err != nil {
			// 可选：忽略单个路径错误，继续遍历其他文件（根据需求调整）
			fmt.Printf("警告：遍历路径 %s 失败，跳过: %v\n", path, err)
			return nil
		}

		// 只处理文件（跳过目录）
		if d.IsDir() {
			return nil
		}

		// 4. 转换为绝对路径（关键步骤）
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Printf("警告：无法获取 %s 的绝对路径，跳过: %v\n", path, err)
			return nil
		}

		// 5. 添加到结果列表
		fileAbsPaths = append(fileAbsPaths, absPath)
		return nil
	})

	// 6. 处理遍历的全局错误（如整体权限拒绝）
	if err != nil {
		return nil, fmt.Errorf("遍历目录失败: %w", err)
	}

	return fileAbsPaths, nil
}
