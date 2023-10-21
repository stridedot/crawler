package collect

import (
	"bufio"
	"fmt"
	"github.com/stridedot/crawler/proxy"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
	"time"
)

type Fetcher interface {
	Get(url string) ([]byte, error)
}

type BaseFetch struct{}

type BrowserFetch struct {
	Timeout time.Duration
	Proxy   proxy.Func
}

// Get 爬取数据
func (f *BaseFetch) Get(url string) ([]byte, error) {
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

// Get 通过创建客户端的方式爬取数据
func (f BrowserFetch) Get(url string) ([]byte, error) {
	client := http.Client{
		Timeout: f.Timeout,
	}

	if f.Proxy != nil {
		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = f.Proxy
		client.Transport = transport
	}
	
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Fetch url error: %v\n", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Fetch url error: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error status code: %v\n", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())

	return io.ReadAll(utf8Reader)
}

// DetermineEncoding 判断爬取数据的编码
func DetermineEncoding(r *bufio.Reader) encoding.Encoding {
	b, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(b, "")
	return e
}
