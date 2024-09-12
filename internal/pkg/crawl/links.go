package crawl

import (
	"io"
	"log"
	"netwatch/internal/pkg/config"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func LinkExtractor(
	body io.ReadCloser,
	roam bool,
	linkConfig *config.LinksConfig) ([]string, error) {
	if !linkConfig.Crawl { // TODO: this is just the per-site config, check global too
		return []string{}, nil
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
		return []string{}, err
	}

	var links []string
	doc.Find("a[href]").Each(func(index int, element *goquery.Selection) {
		href, exists := element.Attr("href")
		if exists && (roam || (linkFilter(href) && linkPatternMatch(href, *linkConfig))) {
			links = append(links, href)
		}
	})

	return links, nil
}

// Removes links that we don't want to crawl
func linkFilter(link string) bool {
	if strings.HasPrefix(link, "http://") ||
		strings.HasPrefix(link, "https://") ||
		strings.HasPrefix(link, "/") {
		return true
	}

	return false
}

// Returns true for links that match the pattern provided in the config
func linkPatternMatch(link string, linkConfig config.LinksConfig) bool {
	if linkConfig.Pattern == "" || linkConfig.Pattern == "*" {
		return true
	}

	regex, err := regexp.Compile(linkConfig.Pattern)
	if err != nil {
		log.Println("Error compiling link pattern")
		return false
	}

	return regex.MatchString(link)
}
