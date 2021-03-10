package house

import (
	"HouseSpider/conf"
	"HouseSpider/proxy"
	"HouseSpider/request"
	"errors"
	"fmt"
	"log"
	"time"
)

type House struct {
	validUrls []string
	conf      *conf.Config
	proxy     *proxy.Context
	request   *request.Request
}

func New(conf *conf.Config, proxy *proxy.Context) *House {
	h := &House{conf: conf, proxy: proxy}
	h.request = request.New(h.httpCallback)
	return h
}

func (h *House) handleOverview(url, body string, err error) {
	if err != nil {
		log.Println(url)
	} else {
		log.Printf("err %s", url)
	}
}

func (h *House) httpCallback(url, body string, err error, opaque interface{}) {
	if opaque.(string) == "overview" {
		h.handleOverview(url, body, err)
	}
}

func (h *House) getProxy() (string, error) {
	for i := 0; i < h.conf.MaxProxyRetry; i++ {
		proxyUrl, err := h.proxy.Get()
		if err == nil {
			return proxyUrl, nil
		}
		log.Printf("there is no proxy avalible, sleep %d seconds", h.conf.WaitProxyTime)
		time.Sleep(time.Duration(h.conf.WaitProxyTime) * time.Second)
	}
	return "", errors.New("get proxy timeout")
}

func (h *House) Fetch() {
	for _, url := range h.conf.StartUrls {
		url += h.conf.Keyword
		start := 0
		for i := 0; i < h.conf.MaxReq; i++ {
			proxyUrl, err := h.getProxy()
			if err != nil {
				continue
			}
			url += "&start=" + fmt.Sprint(start)
			h.request.AsyncGet(url, proxyUrl, "overview")
			start += 50
		}
	}
	h.request.WaitAllDone()
}
