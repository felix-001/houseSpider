package proxy

import (
	"HouseSpider/request"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type RawProxy struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Type      string `json:"type"`
	Anonymity string `json:"anonymity"`
}

type Proxy struct {
	url       string
	timestamp int
	reqErrCnt int
}

type Context struct {
	proxies map[string]Proxy
	url     string
	ip      string
}

func New(url string) *Context {
	return &Context{url: url, proxies: map[string]Proxy{}}
}

func (ctx *Context) parse(html string) {
	reader := bufio.NewReader(strings.NewReader(html))
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Printf("parse json proxies done!")
			break
		}
		var rawProxy RawProxy
		err = json.Unmarshal(line, &rawProxy)
		if err != nil {
			log.Fatal("unmarsha json error")
			continue
		}
		if rawProxy.Type == "https" {
			url := "http://" + rawProxy.Host + ":" + fmt.Sprint(rawProxy.Port)
			proxy := Proxy{url: url}
			ctx.proxies[url] = proxy
		}
	}
	log.Printf("total got %d proxies", len(ctx.proxies))
}

func (ctx *Context) Fetch() {
	req := request.New(nil)
	html, err := req.Get(ctx.url, "")
	if err != nil {
		log.Println(err)
		return
	}
	ctx.parse(html)
}

const ProxyWaitTime = 5

func (ctx *Context) isProxyAvailable(proxy *Proxy) bool {
	if proxy.timestamp == 0 {
		return true
	}
	curTimestamp := time.Now().Second()
	if curTimestamp-proxy.timestamp < ProxyWaitTime {
		return false
	}
	return true
}

var ErrNoAvailableProxy = errors.New("no availble proxy")

func (ctx *Context) Get() (string, error) {
	for _, proxy := range ctx.proxies {
		if !ctx.isProxyAvailable(&proxy) {
			continue
		}
		proxy.timestamp = time.Now().Second()
		return proxy.url, nil
	}
	return "", ErrNoAvailableProxy
}

var maxErrCnt = 5

// IncErrCnt inc err cnt
func (ctx *Context) IncErrCnt(url string) {
	proxy := ctx.proxies[url]
	proxy.reqErrCnt++
	if proxy.reqErrCnt == maxErrCnt {
		log.Fatalf("proxy: %s err cnt reach max, delete it", url)
		delete(ctx.proxies, url)
	}
}

func (ctx *Context) callback(urlStr, body string, err error, opaque interface{}) {
	if err != nil {
		delete(ctx.proxies, opaque.(string))
		return
	}
}

func (ctx *Context) getIp() {
	request := request.New(nil)
	request.Get("")
}

func (ctx *Context) Filter() {
	req := request.New(ctx.callback)
	for _, proxy := range ctx.proxies {
		req.AsyncGet("http://httpbin.org/get", proxy.url, proxy.url)
	}
	req.WaitAllDone()
	log.Printf("after filter, valid proxy count: %d", len(ctx.proxies))
}
