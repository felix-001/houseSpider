package proxy

import (
	"HouseSpider/httpreq"
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
}

func New(url string) *Context {
	return &Context{url: url}
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
		url := rawProxy.Type + "://" + rawProxy.Host + ":" + fmt.Sprint(rawProxy.Port)
		proxy := Proxy{url: url}
		ctx.proxies[url] = proxy
	}
}

func (ctx *Context) Fetch() {
	html, err := httpreq.Get(ctx.url, "")
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
