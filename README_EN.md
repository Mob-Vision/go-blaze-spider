English | [中文](/README.md)

# A lightweight crawl go framework

a lightweight crawl framework written in go


# Install

> go get github.com/bennya8/go-blaze-spider

# Usage

@SEE https://github.com/bennya8/go-blaze-spider-example

```go
package main

import (
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	go_blase_spider "github.com/bennya8/go-blaze-spider"
	"github.com/imroc/req"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var (
	db  *sql.DB
	err error
)

func init() {
	db, err = sql.Open("sqlite3", "db/stackoverflow.db")
	if err != nil {
		log.Fatalln(err)
	}
	createTable()
}

func createTable() {

	sql := `CREATE TABLE "stackoverflow_job" ("id" integer,"logo" varchar,"title" varchar,"firm" varchar,"summary" text, PRIMARY KEY (id))`
	_, err := db.Exec(sql)
	if err != nil {
		log.Println(err)
	}
}
func main() {

	spiderOps := []go_blase_spider.TaskOpt{
		go_blase_spider.TaskOptEnableCookie(true),
		//go_spider.TaskOptSrcCharset("gbk"),
		go_blase_spider.TaskOptGapLimitRandom(500, 5000),
		go_blase_spider.TaskOptGapLimit(1000),
		//go_spider.TaskOptProxy([]string{"127.0.0.1:8700"}),
		go_blase_spider.TaskOptCache("cache"),
		// domain white list
		go_blase_spider.TaskOptDomains([]string{"stackoverflow.com", "www.stackoverflow.com"}),
	}

	// spider entry point
	task := go_blase_spider.NewTaskHandler("stackoverflow", "https://stackoverflow.com/jobs", spiderOps...)

	// setting up fake ua
	//task.Headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	//task.Headers["Cookie"] = `uu=eyJpZCI6InV1N2JlNzkwMDVkYmMyNGYwOTk3ZDMiLCJwcmVmZXJlbmNlcyI6eyJmaW5kX2luY2x1ZGVfYWR1bHQiOmZhbHNlfX0=; adblk=adblk_no; session-id=131-6048916-9533330; session-id-time=2252674724; csm-hit=tb:X3AS0B8X80FQBMSJNNAZ+s-X3AS0B8X80FQBMSJNNAZ|1621954740210&t:1621954740210&adb:adblk_no`
	//task.Headers["User-Agent"] = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"

	// callback before each request
	task.OnRequest(func(url string, header req.Header, param req.Param, err error) {
		fmt.Println("OnRequest", url, err)
	})

	// callback after each request
	task.OnResponse(func(url string, resp *req.Resp, err error) {
		fmt.Println("OnResponse", url, err)
	})

	// LOGIC START

	prefix := "https://stackoverflow.com"

	// 1. fetch job list cell.
	task.OnQuery(".listResults", func(url string, selection *goquery.Selection) {

		selection.Each(func(i int, selection *goquery.Selection) {
			selection.Find(".grid").Each(func(i int, selection *goquery.Selection) {

				logo, _ := selection.Find(".w48.h48.bar-sm").Attr("src")
				title := selection.Find(".mb4.fc-black-800.fs-body3 a").Text()
				titleUrl, exist := selection.Find(".mb4.fc-black-800.fs-body3 a").Attr("href")
				if exist {
					task.Visit(prefix + titleUrl)
				}

				fmt.Println(logo)
				fmt.Println(title)

				// writing record to table with sqlite.
				stmt, err := db.Prepare("INSERT INTO stackoverflow_job(logo, title) values (?,?)")
				if err != nil {
					fmt.Println(err)
				}
				rs, err := stmt.Exec(logo, title)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(rs)
			})
		})
	})

	// 1.2 simulate click next button
	task.OnQuery(".s-pagination a.s-pagination--item", func(url string, selection *goquery.Selection) {
		last := selection.Last()
		nextUrl, exist := last.Attr("href")
		if exist {
			task.Visit(prefix + nextUrl)
		}
	})

	// 2. fetch job detail

	// LOGIC END

	// create main spider
	spider := go_blase_spider.NewGoSpider()

	// adding crawl task to main spider
	spider.AddTask(task)

	// execution
	spider.Run()

}

```

![alt 属性文本](https://github.com/bennya8/go-blaze-spider-example/blob/master/screenshot/WX20210525-233428.png)

# Change log

## v1.0.2

* [IMPROVE] request random timegap

## v1.0.1

* [ADDED] url cache feature
* [ADDED] page encode convert feature

## v1.0.0.alpha

* build framework skeleton
* [TODO] mock ua/ url caching

# 3rd dependencies

- github.com/imroc/req -- an effective go http request library
- github.com/PuerkitoBio/goquery - dom parser
- github.com/axgle/mahonia - charset converter

