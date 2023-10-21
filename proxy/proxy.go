package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

type Func func(*http.Request) (*url.URL, error)

type roundRobinSwitcher struct {
	proxyURLs []*url.URL
	index     uint32
}

func (rr *roundRobinSwitcher) GetProxy(*http.Request) (*url.URL, error) {
	// atomic.AddUint32 函数安全地增加 `rr.index` 的值
	// 该函数是个原子操作，可以在多个 goroutine 中安全地使用
	// 增加后的值减去 1，是为了获得当前 index 的值，因为原子增加返回的
	// 是增加后的值
	index := atomic.AddUint32(&rr.index, 1) - 1
	// 通过取模的方式，确保 `index` 的值不会超过 `rr.proxyURLs` 的长度
	// 当 `index` 的值超过 `rr.proxyURLs` 的长度时，会从头开始，
	// 也就是实现了轮询的效果
	u := rr.proxyURLs[index%uint32(len(rr.proxyURLs))]
	return u, nil
}

func RoundRobinProxySwitcher(ProxyURLs ...string) (Func, error) {
	if len(ProxyURLs) == 0 {
		return nil, errors.New("ProxyURLs is empty")
	}

	urls := make([]*url.URL, len(ProxyURLs))
	for i, v := range ProxyURLs {
		parsed, err := url.Parse(v)
		if err != nil {
			return nil, err
		}
		urls[i] = parsed
	}

	rr := &roundRobinSwitcher{urls, 0}

	return rr.GetProxy, nil
}
