package house

import (
	"HouseSpider/conf"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type House struct {
	validUrls []string
	conf      *conf.Config
	c         *colly.Collector
	cnt       int
}

func (h *House) isValid(title string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(title, keyword) {
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
		if !h.isValid(e.Text, h.conf.FilterKeywords) {
			log.Printf("drop %s %s\n", e.Text, link)
			return
		}
		fmt.Printf("%03d %q -> %s\n", h.cnt, e.Text, link)
		if h.cnt < 10 {
			c.Visit(e.Request.AbsoluteURL(link))
			time.Sleep(5 * time.Second)
		}
		h.cnt++
	})

	c.OnHTML("div[class*=rich-content]>p", func(e *colly.HTMLElement) {
		fmt.Printf("%q\n", e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		//r.Headers.Set("Referer", "https://www.douban.com")
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(resp *colly.Response, err error) {
		fmt.Println(err)
	})
	h.c = c
	return h
}

func (h *House) Fetch() {
	h.c.Visit("https://www.douban.com/group/search?cat=1013&group=26926&sort=time&q=%E5%9B%9E%E9%BE%99%E8%A7%82")
}
