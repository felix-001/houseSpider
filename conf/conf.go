package conf

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	StartUrls      []string `json:"start_urls"`
	Keywords       []string `json:"keywords"`
	FilterKeywords []string `json:"filter_keywords"`
	Groups         []string `json:"groups"`
	FilterAuthors  []string `json:"filter_authors"` // 过滤作者id
	MininumChars   int      `json:"mininum_chars"`  // 介绍详情至少要多少个文字
	MaxPages       int      `json:"max_pages"`      // 每个小组最多爬多少个page
	MaxDays        int      `json:"max_days"`       // 过了多少天的帖子就不要了
	BaseUrl        string   `json:"base_url"`
	ReqInterval    int      `json:"req_interval"`
}

func Load(confPath string) (*Config, error) {
	jsonFile, err := os.Open(confPath)
	if err != nil {
		log.Println("Error opening json file:", err)
		return nil, err
	}

	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)
	var config Config
	for {
		err := decoder.Decode(&config)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error decoding json:", err)
			return nil, err
		}

	}
	return &config, err
}
