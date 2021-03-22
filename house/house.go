package house

import (
	"HouseSpider/conf"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Post struct {
	url   string
	title string
}

type House struct {
	validPosts          []Post
	conf                *conf.Config
	c                   *colly.Collector
	totalCnt            int
	contentInvalidCnt   int
	titleInvalidCnt     int
	contentTooLittleCnt int            // 介绍详情文字太少了
	statistics          map[string]int // 统计各个关键字过滤的url个数
	postTooOldCnt       int
	invalidAuthorCnt    int
	originTitle         string
}

func (h *House) isValid(txt string, keywords []string) (bool, string) {
	for _, keyword := range keywords {
		if strings.Contains(txt, keyword) {
			return false, keyword
		}
	}
	return true, ""
}

func (h *House) isKeywordValid(txt string) (bool, string) {
	return h.isValid(txt, h.conf.FilterKeywords)
}

func (h *House) isAuthorValid(author string) (bool, string) {
	return h.isValid(author, h.conf.FilterAuthors)
}

func (s *House) str2time(timestr string) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", timestr, loc)
	return t, err
}

func New(conf *conf.Config) *House {
	h := &House{conf: conf, statistics: map[string]int{}}
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36"),
	)
	c.OnHTML("tr[class=pl]", func(e *colly.HTMLElement) {
		defer func() {
			h.totalCnt++
		}()
		link := e.ChildAttr("a[class]", "href")
		title := e.ChildAttr("td[class=td-subject]>a", "title")
		h.originTitle = title
		if valid, keyword := h.isKeywordValid(title); !valid {
			h.statistics[keyword]++
			log.Printf("drop %s, invalid keyword found in search page, %s, keyword: %s\n", link, title, keyword)
			return
		}

		timeStr := e.ChildAttr("td[class=td-time]", "title")
		tm, _ := h.str2time(timeStr)
		now := time.Now()
		duration := now.Sub(tm)
		if duration.Hours() > float64(h.conf.MaxDays)*24 {
			log.Printf("drop %s, post too old, time: %s\n", link, timeStr)
			h.postTooOldCnt++
			return
		}
		c.Visit(e.Request.AbsoluteURL(link))
		time.Sleep(time.Duration(h.conf.ReqInterval) * time.Second)
	})

	c.OnHTML("div[class=topic-doc]", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		author := e.ChildAttr("span[class=from]>a", "href")
		if valid, _author := h.isAuthorValid(author); !valid {
			log.Printf("drop %s, invalid author found, %s, author: %s", url, author, _author)
			h.invalidAuthorCnt++
			return
		}
		title := e.ChildText("td[class=tablecc]")
		if valid, keyword := h.isKeywordValid(title); !valid {
			log.Printf("drop %s, invalid keyword found in title, %s, keyword: %s", url, title, keyword)
			h.statistics[keyword]++
			h.titleInvalidCnt++
			return
		}
		content := e.ChildText("div[class=topic-content]")
		if len(content) < h.conf.MininumChars {
			log.Printf("drop %s, content len < %d actual: %d, content: %s", url, h.conf.MininumChars, len(content), content)
			h.contentTooLittleCnt++
			return
		}
		if valid, keyword := h.isKeywordValid(content); !valid {
			log.Printf("drop %s, invalid keyword found in content detail, keyword:%s, content: %s", url, keyword, content)
			h.statistics[keyword]++
			h.contentInvalidCnt++
			return
		}
		if title == "" {
			// 有的详情的帖子标题不是放到topic-doc里面的
			// 是放在div[class=otice-info]>h1里面的
			title = h.originTitle
		}
		post := Post{url: url, title: title}
		h.validPosts = append(h.validPosts, post)
		log.Printf("append url: %s to valid url list", url)
	})

	c.OnRequest(func(r *colly.Request) {
		url, _ := url.QueryUnescape(r.URL.String())
		log.Printf("+++++ [%04d] Visiting %s\n", h.totalCnt, url)
	})

	c.OnError(func(resp *colly.Response, err error) {
		log.Println(err)
	})
	h.c = c
	return h
}

func (h *House) saveHtml() {
	str := ""
	for i, post := range h.validPosts {
		str += "<a href=\"" + post.url + "\">" + "[" + fmt.Sprint(i) + "]" + post.title + "</a><br>\n"
	}
	if err := ioutil.WriteFile("./houses.html", []byte(str), 0666); err != nil {
		log.Println(err)
	}
}

func (h *House) dumpStatistics() {
	log.Printf("content chars too little count: %d\n", h.contentTooLittleCnt)
	log.Printf("content invalid count: %d\n", h.contentInvalidCnt)
	log.Printf("title invalid count: %d\n", h.titleInvalidCnt)
	log.Printf("post too old count: %d\n", h.postTooOldCnt)
	log.Printf("invalid author count: %d\n", h.invalidAuthorCnt)
	log.Printf("total: %d\n", h.totalCnt)
	for k, v := range h.statistics {
		log.Printf("%s ==> %d\n", k, v)
	}
}

func (h *House) genQuerys() []string {
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
	return querys
}

func (h *House) Fetch() {
	querys := h.genQuerys()
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
				time.Sleep(time.Duration(h.conf.ReqInterval) * time.Second)
			}
		}
	}
	h.dumpStatistics()
	h.saveHtml()
}
