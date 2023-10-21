package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/stridedot/crawler/collect"
	"github.com/stridedot/crawler/proxy"
	"time"
)

func main() {
	proxyURLs := []string{"http://127.0.0.1:7890", "http://127.0.0.1:8889"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		panic(err)
	}
	url := "https://www.baidu.com"
	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}
	body, err := f.Get(url)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(body))

	// 加载HTML文档
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("read content failed:%v\n", err)
	}

	doc.Find("div.news_li h2 a[target=_blank]").Each(func(i int, s *goquery.Selection) {
		// 获取匹配元素的文本
		title := s.Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}
