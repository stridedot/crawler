package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/stridedot/crawler/collect"
)

func main() {
	url := "https://news.sina.com.cn/"
	var f collect.Fetcher = &collect.BaseFetch{}

	body, err := f.Get(url)
	if err != nil {
		fmt.Printf("Fetcher error: %v\n", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	doc.Find(`h1[data-client="headline"] a[target="_blank"]`).Each(func(i int, s *goquery.Selection) {
		fmt.Printf("Fetch title: %s\n", s.Text())
	})
}
