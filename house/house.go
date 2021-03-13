package house

import (
	"HouseSpider/conf"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type House struct {
	validUrls           []string
	conf                *conf.Config
	c                   *colly.Collector
	totalCnt            int
	title               string
	content             string
	contentInvalidCnt   int
	titleInvalidCnt     int
	contentTooLittleCnt int            // 介绍详情文字太少了
	statistics          map[string]int // 统计各个关键字过滤的url个数
}

func (h *House) isValid(txt string) (bool, string) {
	for _, keyword := range h.conf.FilterKeywords {
		if strings.Contains(txt, keyword) {
			return false, keyword
		}
	}
	return true, ""
}

func (s *House) str2time(timestr string) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, err := time.ParseInLocation("2006-01-02T15:04:05", timestr, loc)
	return t, err
}

func New(conf *conf.Config) *House {
	h := &House{conf: conf, statistics: map[string]int{}}
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36"),
	)
	c.OnHTML("tr[class=pl]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.ChildAttr("td[class=td-subject]>a", "title")
		if valid, keyword := h.isValid(title); !valid {
			h.statistics[keyword]++
			log.Printf("drop %s invalid keyword found in search page, %s, keyword:%s\n", link, title, keyword)
			return
		}
		timeStr := e.ChildAttr("td[class=td-time]", "title")
		tm, _ := h.str2time(timeStr)
		now := time.Now()
		duration := now.Sub(tm)
		if duration.Hours() > 7*24 {
			log.Printf("drop %s, post too old, time: %s\n", link, timeStr)
			return
		}
		c.Visit(e.Request.AbsoluteURL(link))
		time.Sleep(3 * time.Second)
	})

	c.OnHTML("div[class*=rich-content]>p", func(e *colly.HTMLElement) {
		h.content = e.Text
	})

	c.OnHTML("td[class=tablecc]", func(e *colly.HTMLElement) {
		h.title = e.Text
	})

	c.OnRequest(func(r *colly.Request) {
		url, _ := url.QueryUnescape(r.URL.String())
		log.Printf("+++++ [%03d] Visiting %s\n", h.totalCnt, url)
	})

	c.OnScraped(func(resp *colly.Response) {
		url := resp.Request.URL.String()
		defer func() {
			h.title = ""
			h.content = ""
		}()
		if strings.Contains(url, "cat") {
			return
		}
		h.totalCnt++
		if len(h.content) < h.conf.MininumChars {
			log.Printf("drop %s, content len < %d actual: %d, content: %s", url, h.conf.MininumChars, len(h.content), h.content)
			h.contentTooLittleCnt++
			return
		}
		if valid, keyword := h.isValid(h.content); !valid {
			log.Printf("drop %s, invalid keyword found in content detail, keyword:%s, content: %s", url, keyword, h.content)
			h.statistics[keyword]++
			h.contentInvalidCnt++
			return
		}
		if valid, keyword := h.isValid(h.title); !valid {
			log.Printf("drop %s, invalid keyword found in title, %s, keyword: %s", url, h.title, keyword)
			h.statistics[keyword]++
			h.titleInvalidCnt++
			return
		}
		h.validUrls = append(h.validUrls, url)
	})

	c.OnError(func(resp *colly.Response, err error) {
		log.Println(err)
	})
	h.c = c
	return h
}

func (h *House) Fetch() {
	querys := []string{}
	query := ""
	i := 1
	for _, keyword := range h.conf.Keywords {
		query += keyword + " "
		if i == 3 {
			querys = append(querys, query)
			query = ""
			i = 0
		}
		i++
	}
	if query != "" {
		querys = append(querys, query)
	}
	for _, group := range h.conf.Groups {
		for _, query := range querys {
			encodeQuery := url.QueryEscape(query)
			for i := 0; i < h.conf.MaxPages; i++ {
				urlNew := h.conf.BaseUrl +
					"&group=" + group +
					"&q=" + encodeQuery +
					"&start=" + fmt.Sprint(i*50)
				log.Println(urlNew)
				h.c.Visit(urlNew)
				time.Sleep(3 * time.Second)
			}
		}
	}
	//h.c.Visit("https://www.douban.com/group/search?cat=1013&group=26926&sort=time&q=%E5%9B%9E%E9%BE%99%E8%A7%82")
	log.Printf("content chars too little count: %d\n", h.contentTooLittleCnt)
	log.Printf("content invalid count: %d\n", h.contentInvalidCnt)
	log.Printf("title invalid count: %d\n", h.titleInvalidCnt)
	log.Printf("total: %d", h.totalCnt)
	for k, v := range h.statistics {
		log.Printf("%s ==> %d", k, v)
	}
	log.Println(h.validUrls)
}
