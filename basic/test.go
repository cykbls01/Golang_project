package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
)

func main() {
	// 创建一个饼状图实例
	pie := charts.NewPie()

	// 设置全局配置项
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "我的饼状图示例"}), // 设置图表标题
	)

	// 添加系列数据，这些是饼状图的各个部分
	pie.AddSeries("访问来源",
		[]opts.PieData{
			{Value: 1048, Name: "搜索引擎 1048"},
			{Value: 735, Name: "直接访问 735"},
			{Value: 580, Name: "邮件营销 580"},
			{Value: 484, Name: "联盟广告 484"},
			{Value: 300, Name: "视频广告 300"},
		},
	)
	f, _ := os.Create("pie.html")
	// 渲染图表到HTML文件中
	if err := pie.Render(f); err != nil {
		panic(err)
	}
}
