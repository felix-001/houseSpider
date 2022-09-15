package house

import (
	"log"

	"github.com/gocolly/colly"
)

const (
	uri = "https://spys.one/en/free-proxy-list/"
)

type SPYSAgent struct {
}

func (s *SPYSAgent) GetUrl() string {
	return uri
}

func (s *SPYSAgent) GetSelector() string {
	return ".spy14"
}

func (s *SPYSAgent) HTMLCallback(e *colly.HTMLElement) {
	log.Println(e.Name, e.Text)
}
