package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
)

func main() {
	url := "https://www.thepaper.cn/"
	body, err := Fetcher(url)
	if err != nil {
		fmt.Printf("Fetcher error: %v\n", err)
		return
	}

	fmt.Println(string(body))
}

// Fetcher 爬取数据
func Fetcher(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Fetch url error: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error status code: %v\n", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)

	// 爬取数据的编码
	e := DetermineEncoding(bodyReader)

	// 转换编码
	utf8Reader := transform.NewReader(resp.Body, e.NewDecoder())

	return io.ReadAll(utf8Reader)
}

// DetermineEncoding 判断爬取数据的编码
func DetermineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
