package main

import (
	"github.com/stridedot/crawler/collect"
	"github.com/stridedot/crawler/log"
	"github.com/stridedot/crawler/proxy"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func main() {
	plugin, closer := log.NewFilePlugin("./test.log", zapcore.InfoLevel)
	defer closer.Close()
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	proxyURLs := []string{"http://127.0.0.1:7890", "http://127.0.0.1:8889"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed", zap.Error(err))
	}
	url := "https://www.baidu.com"
	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}
	body, err := f.Get(url)
	if err != nil {
		logger.Error("read content failed", zap.Error(err))
	}

	logger.Info("get content", zap.Int("len", len(body)))
}
