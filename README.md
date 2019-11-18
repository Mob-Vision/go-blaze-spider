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

github.com/imroc/req

**effective go http request library**


github.com/PuerkitoBio/goquery

**dom parser**


github.com/axgle/mahonia

**character set converter**

# Disclaimer

All downloads and use of the software (bennya8/go-spider) are deemed to have been read carefully and fully agree to the following disclaimer:

The software (bennya8/go-spider) is for personal learning and communication purposes only and is strictly prohibited for commercial and non-defective purposes.

The author of the software (bennya8/go-spider) assumes no responsibility for any commercial activity or non-use, and all risks will be borne by him or her.

The software (bennya8/go-spider), except for the terms of service, any accidents, negligence, contract damage, defects, copyright or other intellectual property infringement caused by improper use of the software, and any damage caused by the software, the software The author is not responsible and does not assume any legal responsibility.

For issues not covered by this statement, please refer to relevant national laws and regulations. When this statement conflicts with relevant national laws and regulations, the national laws and regulations shall prevail.

Copyright of this software and its right of modification, renewal and final interpretation are owned by the author of the software (bennya8/go-spider).



