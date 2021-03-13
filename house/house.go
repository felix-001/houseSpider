package house

import (
	"HouseSpider/conf"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type House struct {
	validUrls           []string
	conf                *conf.Config
	c                   *colly.Collector
	cnt                 int
	title               string
	content             string
	contentInvalidCnt   int
	titleInvalidCnt     int
	contentTooLittleCnt int // 介绍详情文字太少了
}

func (h *House) isValid(txt string) bool {
	for _, keyword := range h.conf.FilterKeywords {
		if strings.Contains(txt, keyword) {
			return false
		}
	}
	return true
}

func New(conf *conf.Config) *House {
	h := &House{conf: conf}
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36"),
	)
	c.OnHTML("td[class=td-subject]>a[class]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if !h.isValid(e.Text) {
			log.Printf("drop %s %s\n", e.Text, link)
			return
		}
		log.Printf("%03d %q -> %s\n", h.cnt, e.Text, link)
		if h.cnt < 10 {
			c.Visit(e.Request.AbsoluteURL(link))
			time.Sleep(5 * time.Second)
		}
		h.cnt++
	})

	c.OnHTML("div[class*=rich-content]>p", func(e *colly.HTMLElement) {
		log.Printf("%q\n", e.Text)
		h.content = e.Text
	})

	c.OnHTML("td[class=tablecc]", func(e *colly.HTMLElement) {
		log.Printf("%s\n", e.Text)
		h.title = e.Text
		//h.saveUrl(e.Text, e.Request.URL.String())
	})

	c.OnRequest(func(r *colly.Request) {
		//r.Headers.Set("Referer", "https://www.douban.com")
		log.Println("Visiting", r.URL)
	})

	c.OnScraped(func(resp *colly.Response) {
		url := resp.Request.URL.String()
		if strings.Contains(url, "cat") {
			return
		}
		if len(h.content) < h.conf.MininumChars {
			h.contentTooLittleCnt++
			return
		}
		if !h.isValid(h.content) {
			h.contentInvalidCnt++
			return
		}
		if !h.isValid(h.title) {
			h.titleInvalidCnt++
			return
		}
		h.validUrls = append(h.validUrls, url)
		h.title = ""
		h.content = ""
	})

	c.OnError(func(resp *colly.Response, err error) {
		log.Println(err)
	})
	h.c = c
	return h
}

func (h *House) Fetch() {
	h.c.Visit("https://www.douban.com/group/search?cat=1013&group=26926&sort=time&q=%E5%9B%9E%E9%BE%99%E8%A7%82")
	log.Println(h.validUrls)
	log.Printf("content chars too little count: %d\n", h.contentTooLittleCnt)
	log.Printf("content invalid count: %d\n", h.contentInvalidCnt)
	log.Printf("title invalid count: %d\n", h.titleInvalidCnt)
}
