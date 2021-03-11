package house

import (
	"HouseSpider/conf"
	"HouseSpider/proxy"
	"HouseSpider/request"
	"bytes"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
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
		log.Println(err)
		return
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(body)))
	if err != nil {
		log.Println(err)
		return
	}
	doc.Find("tr[class=pl]").Each(func(i int, selection *goquery.Selection) {
		//log.Println(selection.Text())
	})
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
		start := 0
		for i := 0; i < h.conf.MaxReq; i++ {
			proxyUrl, err := h.getProxy()
			if err != nil {
				continue
			}
			urlNew := url + h.conf.Keyword + "&start=" + fmt.Sprint(start)
			log.Println(proxyUrl)
			h.request.AsyncGet(urlNew, proxyUrl, "overview")
			start += 50
		}
	}
	h.request.WaitAllDone()
}
