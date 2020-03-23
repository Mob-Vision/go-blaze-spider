English | [中文](/README.md)

# A lightweight crawl go framework

a lightweight crawl framework written in go


# Install

> go get github.com/bennya8/go-spider

# Usage

```golang
    // Configuare crawl task options 
    job51Opts := []TaskOpt{
        TaskOptEnableCookie(true),
        TaskOptGapLimit(5000),
        TaskOptCache("cache"),
        TaskOptProxy([]string{"127.0.0.1:8700"}),
        TaskOptSrcCharset("gbk"),
        TaskOptDomains([]string{"www.51job.com", "search.51job.com", "jobs.51job.com"}),
    }
    // Craete new task handler and passing options，NewTaskHandler(name string,entry string,opts ...opts）
    job51 := NewTaskHandler("job51", "https://www.51job.com", job51Opts...)

    // Before request event
    job51.OnRequest(func(url string, header *req.Header, param *req.Param, err error) {
        fmt.Println(url, header, param, err)
    })

    // After request event
    job51.OnResponse(func(resp *req.Resp, err error) {
        fmt.Println(resp, err)
    })
    
    // Dom search 
    // allowing nest selection, check （github.com/PuerkitoBio/goquery）to get more example
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
    
    // create main spider thread
    spider := NewGoSpider()
    
    // register current task to the main spider thread
    // supported muti-tasking
    spider.AddTask(job51)

    // execution 
    spider.Run()
```


# Change log

## v1.0.1

* [ADDED] url cache feature
* [ADDED] page encode convert feature

## v1.0.0.alpha (2019-10-05 22:22 UTC+8:00)

* build framework skeleton
* [TODO] mock ua/ url caching

# 3rd dependencies

- github.com/imroc/req -- an effective go http request library
- github.com/PuerkitoBio/goquery - dom parser
- github.com/axgle/mahonia - charset converter

