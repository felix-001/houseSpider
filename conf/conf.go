package conf

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	ProxyUrl       string   `json:"proxy_url"`
	StartUrls      []string `json:"start_urls"`
	MaxReq         int      `json:"max_req"`
	Keyword        string   `json:"keywords"`
	FilterKeywords []string `json:"filter_keywords"`
	MaxProxyRetry  int      `json:"max_proxy_retry"`
	WaitProxyTime  int      `json:"wait_proxy_time"`
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
