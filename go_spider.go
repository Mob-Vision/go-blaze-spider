package go_spider

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/imroc/req"
	net_url "net/url"
	"strings"
	"sync"
	"time"
)

var instance *GoSpider
var once sync.Once

// task handler
type TaskHandler struct {
	Name            string
	Entry           string
	GapLimit        int
	IdleLimit       int
	WorkerLimit     int
	SrcCharset      string
	reqCb           OnRequestCallback
	rspCb           OnResponseCallback
	queryCbs        map[string]OnQueryCallback
	Http            *req.Req
	Headers         *req.Header
	Params          *req.Param
	Domains         []string
	Queue           chan string
	QueueProcessNum int
	QueueTotalNum   int
	Cache           string
}

func (t *TaskHandler) OnRequest(cb OnRequestCallback) {
	t.reqCb = cb
}

func (t *TaskHandler) OnResponse(cb OnResponseCallback) {
	t.rspCb = cb
}

func (t *TaskHandler) OnQuery(selector string, cb OnQueryCallback) {
	t.queryCbs[selector] = cb
}

func (t *TaskHandler) Handle() {
	for {
		select {
		case v := <-t.Queue:
			t.request(v)
			t.QueueProcessNum++
			fmt.Println("total: ", t.QueueTotalNum, " process:", t.QueueProcessNum)
			time.Sleep(time.Millisecond * time.Duration(t.GapLimit))
		default:
			time.Sleep(time.Millisecond * time.Duration(t.IdleLimit))
		}
	}
}

func (t *TaskHandler) request(url string) {

	// fixed abnormal url protocol
	if strings.HasPrefix(url, "//") {
		url = "http://" + strings.Replace(url, "//", "", -1)
	}

	u, err := net_url.Parse(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	domain := u.Hostname()
	isAllow := false
	for _, d := range t.Domains {
		if d == domain {
			isAllow = true
			break
		}
	}
	if !isAllow {
		fmt.Println("url domain not allow", url)
		return
	}

	t.reqCb(url, t.Headers, t.Params, nil)

	resp, err := t.Http.Get(url, t.Headers, t.Params)

	t.rspCb(url, resp, err)

	if err != nil {
		fmt.Println(err)
		return
	}

	var content string
	if len(t.SrcCharset) > 0 {
		content = gbk2utf8(string(resp.String()), t.SrcCharset, "utf-8")
	} else {
		content = string(resp.String())
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(content))
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, q := range t.queryCbs {
		q(url, doc.Find(k))
	}
}

func (t *TaskHandler) Clone() *TaskHandler {
	clone := t
	return clone
}

func (t *TaskHandler) Visit(url string) {
	t.QueueTotalNum++
	go func() {
		t.Queue <- url
	}()
}

func TaskOptSrcCharset(charset string) TaskOpt {
	return func(handler *TaskHandler) {
		if len(charset) > 0 {
			handler.SrcCharset = charset
			return
		}
	}
}

func TaskOptGapLimit(num int) TaskOpt {
	return func(handler *TaskHandler) {
		if num > 0 {
			handler.GapLimit = num
			return
		}
	}
}

func TaskOptEnableCookie(b bool) TaskOpt {
	return func(handler *TaskHandler) {
		if b {
			handler.Http.EnableCookie(true)
		}
	}
}

func TaskOptDomains(domains []string) TaskOpt {
	return func(handler *TaskHandler) {
		if len(domains) > 0 {
			for _, v := range domains {
				handler.Domains = append(handler.Domains, v)
			}
		}
	}
}

type TaskOpt func(*TaskHandler)

func NewTaskHandler(name string, entry string, opts ...TaskOpt) *TaskHandler {

	taskHandler := &TaskHandler{
		Name:            name,
		Entry:           entry,
		Http:            req.New(),
		Headers:         &req.Header{},
		Params:          &req.Param{},
		Queue:           make(chan string),
		QueueTotalNum:   0,
		QueueProcessNum: 0,
		queryCbs:        make(map[string]OnQueryCallback),
	}

	// setting opts
	for _, opt := range opts {
		opt(taskHandler)
	}
	// reduce cpu resource
	if taskHandler.IdleLimit == 0 {
		taskHandler.IdleLimit = 1000
	}
	// prevent request too frequency will cause ip ban
	if taskHandler.GapLimit == 0 {
		taskHandler.GapLimit = 500
	}

	return taskHandler
}

type OnRequestCallback func(url string, header *req.Header, param *req.Param, err error)
type OnResponseCallback func(url string, rsp *req.Resp, err error)
type OnQueryCallback func(url string, selection *goquery.Selection)
type OnErrorCallback func(msg string, err error)

type GoSpider struct {
	handlers []*TaskHandler
}

func NewGoSpider() *GoSpider {
	once.Do(func() {
		instance = new(GoSpider)
		instance.handlers = []*TaskHandler{}
	})
	return instance
}

func (g *GoSpider) AddTask(handler *TaskHandler) {
	g.handlers = append(g.handlers, handler)
}

func (g *GoSpider) Run() {
	for _, v := range g.handlers {
		if len(v.Entry) > 0 {
			go v.Visit(v.Entry)
		}
		go v.Handle()
	}
	select {}
}

func gbk2utf8(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	return string(cdata)
}
