package main

import (
	"fmt"
	"github.com/stridedot/crawler/collect"
	"time"
)

func main() {
	url := "https://book.douban.com/subject/1007305/"
	var f collect.Fetcher = &collect.BrowserFetch{
		Timeout: 3 * time.Second,
	}

	body, err := f.Get(url)
	if err != nil {
		fmt.Printf("Fetcher error: %v\n", err)
		return
	}

	fmt.Println(string(body))

	//doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	//doc.Find(`h1[data-client="headline"] a[target="_blank"]`).Each(func(i int, s *goquery.Selection) {
	//	fmt.Printf("Fetch title: %s\n", s.Text())
	//})
}
