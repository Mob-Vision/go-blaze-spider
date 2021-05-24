中文 | [English](/README_EN.md)

# 轻量爬虫框架

一个用 Golang 实现的轻量级爬虫框架

# 安装

> go get github.com/bennya8/go-blaze-spider

# 使用说明

```golang
    // 配置爬虫任务选项
    job51Opts := []TaskOpt{
        TaskOptEnableCookie(true),
        TaskOptGapLimit(5000),
        TaskOptCache("cache"),
        TaskOptProxy([]string{"127.0.0.1:8700"}),
        TaskOptSrcCharset("gbk"),
        TaskOptDomains([]string{"www.51job.com", "search.51job.com", "jobs.51job.com"}),
    }
    // 创建爬虫任务，NewTaskHandler（命名，入口，选项配置）
    job51 := NewTaskHandler("job51", "https://www.51job.com", job51Opts...)

    // 发起链接前回调
    job51.OnRequest(func(url string, header *req.Header, param *req.Param, err error) {
        fmt.Println(url, header, param, err)
    })

    // 发起链接后回调
    job51.OnResponse(func(resp *req.Resp, err error) {
        fmt.Println(resp, err)
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
    
    // 创建蜘蛛主线程
    spider := NewGoSpider()
    
    // 注册任务
    spider.AddTask(job51)

    // 执行
    spider.Run()
```


# 更新日志


## v1.0.2

* [IMPROVE] 请求访问随机间隔

## v1.0.1

* [ADDED] URL缓存
* [ADDED] 页面编码转换

## v1.0.0.alpha

* 框架锥形骨架
* [TODO] URL缓存/UA设置

# 第三方依赖

- github.com/imroc/req -- 高效的HTTP request库
- github.com/PuerkitoBio/goquery - DOM检索解析器
- github.com/axgle/mahonia - gbk转换utf8 字符编码器

