package collect

import (
	"bufio"
	"fmt"
	"github.com/stridedot/crawler/proxy"
	"go.uber.org/zap"
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
	Get(request *Request) ([]byte, error)
}

type BaseFetch struct{}

type BrowserFetch struct {
	Timeout time.Duration
	Logger  *zap.Logger
	Proxy   proxy.Func
}

// Get 爬取数据
func (f *BaseFetch) Get(request *Request) ([]byte, error) {
	resp, err := http.Get(request.Url)
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
func (f BrowserFetch) Get(request *Request) ([]byte, error) {
	client := http.Client{
		Timeout: f.Timeout,
	}

	if f.Proxy != nil {
		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = f.Proxy
		client.Transport = transport
	}

	req, err := http.NewRequest(http.MethodGet, request.Url, nil)
	if err != nil {
		return nil, fmt.Errorf("Fetch url error: %v\n", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", request.Cookie)

	resp, err := client.Do(req)

	time.Sleep(request.WaitTime)

	if err != nil {
		f.Logger.Error("fetch failed", zap.Error(err))
		return nil, err
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
