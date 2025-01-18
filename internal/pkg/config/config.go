package config

import (
	"errors"
	"log"
	"os"
	"regexp"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Name   string       `yaml:"name"`
	Config GlobalConfig `yaml:"config"`
	Sites  []SiteConfig `yaml:"sites"`
}

type GlobalConfig struct {
	Requests RequestsConfig `yaml:"requests"`
	Recrawl  RecrawlConfig  `yaml:"recrawl"`
	Links    LinksConfig    `yaml:"links"`
	Roam     bool           `yaml:"roam"`
	Queue    QueueConfig    `yaml:"queue"`
}

type RequestsConfig struct {
	Window        int `yaml:"window"`
	MaxConcurrent int `yaml:"maxConcurrent"`
	MaxTotal      int `yaml:"maxTotal"`
	MaxPerHost    int `yaml:"maxPerHost"`
	Timeout       int `yaml:"timeout"`
}

type RecrawlConfig struct {
	Enabled bool `yaml:"enabled"`
	Timeout int  `yaml:"timeout"`
}

type LinksConfig struct {
	Crawl    bool   `yaml:"crawl"`
	Pattern  string `yaml:"pattern"`
	Selector string `yaml:"selector"`
	MaxDepth int    `yaml:"maxDepth"`
}

type QueueConfig struct {
	Capacity int `yaml:"capacity"`
}

type SiteConfig struct {
	URL       string          `yaml:"url"`
	URLRegExp string          `yaml:"urlRegExp"`
	Links     LinksConfig     `yaml:"links"`
	Content   []ContentConfig `yaml:"content"`
}

type ContentConfig struct {
	Name     string `yaml:"name"`
	Selector string `yaml:"selector"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SiteConfigForURL(url string, config *Config) (*SiteConfig, error) {
	for _, site := range config.Sites {
		if site.URLRegExp != "" {
			regex, err := regexp.Compile(site.URLRegExp)
			if err != nil {
				log.Println("Error compiling site URL regex")
				return nil, err
			}

			if regex.MatchString(url) {
				return &site, nil
			}
		} else {
			if site.URL == url {
				return &site, nil
			}
		}
	}

	return nil, errors.New("no site config found for url")
}
