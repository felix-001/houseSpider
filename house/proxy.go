package house

import (
	"log"

	"github.com/gocolly/colly"
)

type Agent interface {
	GetUrl() string
	GetSelector() string
	HTMLCallback(*colly.HTMLElement)
}

type ProxyManager struct {
	agents []Agent
	c      *colly.Collector
}

func NewProxyManager() *ProxyManager {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.1 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/539.38"),
	)
	return &ProxyManager{c: c}
}

func (m *ProxyManager) Register(agent Agent) {
	m.agents = append(m.agents, agent)
}

func (m *ProxyManager) Init() {
	for _, agent := range m.agents {
		m.c.OnHTML(agent.GetSelector(), agent.HTMLCallback)
	}
}

func (m *ProxyManager) Run() error {
	for _, agent := range m.agents {
		url := agent.GetUrl()
		if err := m.c.Visit(url); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
