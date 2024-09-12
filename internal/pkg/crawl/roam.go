package crawl

import (
	"io"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func Roam(body io.ReadCloser) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
		return []string{}, err
	}

	var links []string
	doc.Find("a[href]").Each(func(index int, element *goquery.Selection) {
		href, exists := element.Attr("href")
		if exists {
			links = append(links, href)
		}
	})

	// filter links - ignore #links, mailto:, javascript:, etc
	// match pattern provided in config for this url (if not roam)

	return links, nil
}
