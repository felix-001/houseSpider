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
	valid     bool
}

type Context struct {
	proxies map[string]*Proxy
	url     string
	ip      string
}

func New(url string) *Context {
	return &Context{url: url, proxies: map[string]*Proxy{}}
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
			log.Fatal("unmarsha json error ", err, " ", line)
			continue
		}
		if rawProxy.Type == "https" {
			url := "http://" + rawProxy.Host + ":" + fmt.Sprint(rawProxy.Port)
			proxy := &Proxy{url: url, valid: true}
			ctx.proxies[url] = proxy
		}
	}
	log.Printf("total got %d proxies", len(ctx.proxies))
}

func (ctx *Context) Fetch() {
	req := request.New(nil)
	html, err := req.Get(ctx.url, "", nil)
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
	if !proxy.valid {
		return false
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
		if !ctx.isProxyAvailable(proxy) {
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
		log.Printf("proxy: %s err cnt reach max, delete it", url)
		proxy := ctx.proxies[url]
		proxy.valid = false
	}
}

func (ctx *Context) callback(urlStr, body string, err error, opaque interface{}) {
	if err != nil {
		proxyUrl := opaque.(string)
		log.Println(proxyUrl)
		proxy := ctx.proxies[proxyUrl]
		proxy.valid = false
		return
	}
	if strings.Contains(body, ctx.ip) {
		proxyUrl := opaque.(string)
		log.Println(proxyUrl)
		proxy := ctx.proxies[proxyUrl]
		proxy.valid = false
		return
	}
}

type HTTPBin struct {
	Origin string `json:"origin"`
}

func (ctx *Context) getIP() (string, error) {
	request := request.New(nil)
	resp, err := request.Get("http://httpbin.org/get", "", nil)
	if err != nil {
		log.Println(err)
		return "", err
	}
	httpBin := &HTTPBin{}
	err = json.Unmarshal([]byte(resp), &httpBin)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return httpBin.Origin, nil
}

func (ctx *Context) GetValidProxyLen() int {
	i := 0
	for _, proxy := range ctx.proxies {
		if proxy.valid {
			i++
		}
	}
	return i
}

func (ctx *Context) Filter() {
	ip, err := ctx.getIP()
	if err != nil {
		return
	}
	ctx.ip = ip
	log.Println(ip)
	req := request.New(ctx.callback)
	for _, proxy := range ctx.proxies {
		req.AsyncGet("https://httpbin.org/get", proxy.url, proxy.url)
	}
	req.WaitAllDone()
	log.Printf("after filter, valid proxy count: %d", ctx.GetValidProxyLen())
}
