package crawl

import (
	"io"
	"log"
	"netwatch/internal/pkg/config"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func NewContentExtractor(body io.ReadCloser) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Println(err)
		return doc, err
	}

	return doc, nil
}

func ContentExtractor(doc *goquery.Document, contentConfig *[]config.ContentConfig) (map[string]string, error) {
	content := make(map[string]string)

	for _, contentConfig := range *contentConfig {
		var selector = contentConfig.Selector
		var name = contentConfig.Name

		content[name] = doc.Find(selector).Text()
	}

	return content, nil
}

func LinkExtractor(doc *goquery.Document, url string, roam bool, linkConfig *config.LinksConfig) ([]string, error) {
	var links []string

	if !linkConfig.Crawl { // TODO: this is just the per-site config, check global too - maybe merge global and site config before this point
		return links, nil
	}

	doc.Find("a[href]").Each(func(index int, element *goquery.Selection) {
		href, exists := element.Attr("href")
		if exists && (roam || (linkFilter(href) && linkPatternMatch(href, *linkConfig))) {
			if strings.HasPrefix(href, "/") {
				links = append(links, url+href)
			} else {
				links = append(links, href)
			}
		}
	})

	return links, nil
}

// Removes links that we don't want to crawl
func linkFilter(link string) bool {
	return strings.HasPrefix(link, "http://") ||
		strings.HasPrefix(link, "https://") ||
		strings.HasPrefix(link, "/") ||
		false
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