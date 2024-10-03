package util

import (
	"fmt"
	"time"
)

func JudgeTimeWithTimestamp(timestamp int64, d time.Duration) bool {
	// 将Unix时间戳转换为time.Time类型
	targetTime := time.Unix(timestamp, 0)

	// 获取当前时间
	now := time.Now()

	// 计算当前时间之前一小时的时间点
	compareTime := now.Add(d)

	// 比较时间戳是否在当前时间之前的一小时内
	return targetTime.After(compareTime)
}

func JudgeTimeWithUTC(utcTimeStr string, d time.Duration) bool {
	// 解析UTC时间字符串
	targetTime, err := time.Parse(time.RFC3339, utcTimeStr) // 假设时间字符串符合RFC3339格式，如"2023-04-01T15:04:05Z"
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return false
	}

	now := time.Now()
	// 计算时间差
	compareTime := now.Add(d)

	// 比较时间戳是否在当前时间之前的一小时内
	return targetTime.After(compareTime)
}
