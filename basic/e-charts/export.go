package e_charts

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"math/rand"
	"os"
	"time"
)

// generate random data for bar chart
func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}

func Test() {
	// Step 1: Generate HTML using go-echarts
	htmlFileName := "pie_chart.html"
	if err := generatePieChartHTML(htmlFileName); err != nil {
		fmt.Println("Error generating HTML:", err)
		return
	}

	// Step 2: Convert HTML to PNG using chromedp
	pngFileName := "pie_chart.png"
	if err := convertHTMLToPNG(htmlFileName, pngFileName); err != nil {
		fmt.Println("Error converting HTML to PNG:", err)
		return
	}

	fmt.Println("PNG image generated successfully:", pngFileName)
}
func generatePieChartHTML(filename string) error {
	pie := charts.NewPie()
	pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "Sample Pie Chart",
	}))

	items := []opts.PieData{
		{Name: "Apple", Value: 10},
		{Name: "Banana", Value: 20},
		{Name: "Orange", Value: 30},
		{Name: "Grape", Value: 40},
	}

	pie.SetSeriesOptions(charts.WithPieChartOpts(opts.PieChart{
		Radius: "55%",
	}))

	pie.AddSeries("Category", items)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return pie.Render(f)
}

func convertHTMLToPNG(htmlFileName, pngFileName string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Create a timeout context
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Full path to the HTML file
	htmlFilePath := "file://" + getCurrentDirectory() + "/" + htmlFileName

	var buf []byte
	if err := chromedp.Run(ctx, fullScreenshot(htmlFilePath, 100, &buf)); err != nil {
		return err
	}

	if err := os.WriteFile(pngFileName, buf, 0644); err != nil {
		return err
	}

	return nil
}

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EmulateViewport(900, 500), // 设置浏览器窗口的大小
		chromedp.Navigate(urlstr),
		chromedp.Sleep(2 * time.Second), // wait for the page to fully load
		chromedp.FullScreenshot(res, quality),
	}
}

func getCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}
