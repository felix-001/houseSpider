package conf

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Config struct {
	ProxyUrl string `json:"proxy_url"`
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
			fmt.Println("error decoding json:", err)
			return nil, err
		}

	}
	return &config, err
}
