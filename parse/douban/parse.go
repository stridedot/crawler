package douban

import (
	"github.com/stridedot/crawler/collect"
	"regexp"
)

const detailUrlPattern = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

func ParseURL(content []byte, request *collect.Request) *collect.ParseResult {
	result := &collect.ParseResult{}

	reg := regexp.MustCompile(detailUrlPattern)
	matches := reg.FindAllSubmatch(content, -1)

	for _, m := range matches {
		url := string(m[1])
		result.Requests = append(result.Requests, &collect.Request{
			Url:    url,
			Cookie: request.Cookie,
			ParseFunc: func(c []byte, r *collect.Request) *collect.ParseResult {
				return GetContent(c, url)
			},
		})
	}

	return result
}

const contentPattern = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div>`

func GetContent(content []byte, url string) *collect.ParseResult {
	reg := regexp.MustCompile(contentPattern)
	ok := reg.Match(content)
	if !ok {
		return &collect.ParseResult{
			Items: []interface{}{},
		}
	}

	return &collect.ParseResult{
		Items: []interface{}{url},
	}
}
