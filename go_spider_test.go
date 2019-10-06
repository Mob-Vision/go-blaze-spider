package go_spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
	"testing"
)

func TestSpider(t *testing.T) {
	// 配置爬虫任务选项
	job51Opts := []TaskOpt{
		TaskOptEnableCookie(true),
		TaskOptGapLimit(5000),
		TaskOptDomains([]string{"www.51job.com", "search.51job.com", "jobs.51job.com"}),
	}
	// 创建爬虫任务，NewTaskHandler（命名，入口，选项配置）
	job51 := NewTaskHandler("job51", "https://www.51job.com", job51Opts...)

	// 发起链接前回调
	job51.OnRequest(func(url string, header *req.Header, param *req.Param, err error) {
		fmt.Println(url, header, param, err)
	})

	// 发起链接后回调
	job51.OnResponse(func(url string, resp *req.Resp, err error) {
		//fmt.Println(resp, err)
	})

	// DOM事件注册查询
	//可多个selection嵌套，具体可查看（github.com/PuerkitoBio/goquery）使用方法
	job51.OnQuery(".cn.hlist a", func(url string, selection *goquery.Selection) {
		selection.Each(func(i int, selection *goquery.Selection) {
			href, exists := selection.Attr("href")
			if exists {
				job51.Visit(href)
			}
		})
	})

	job51.OnQuery(".dw_table .el", func(url string, selection *goquery.Selection) {
		selection.Each(func(i int, selection *goquery.Selection) {
			selection.Find("p.t1 a").Each(func(i int, selection *goquery.Selection) {
				href, exists := selection.Attr("href")
				if exists {
					job51.Visit(href)
				}
			})
		})
	})

	//job51.OnQuery("p", func(selection *goquery.Selection) {
	//	selection.Each(func(i int, selection *goquery.Selection) {
	//		fmt.Println(selection.Html())
	//	})
	//})

	// 创建蜘蛛主线程
	spider := NewGoSpider()

	// 注册任务
	spider.AddTask(job51)

	// 执行
	spider.Run()
}
