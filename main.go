package main

import (
	"github.com/stridedot/crawler/collect"
	"github.com/stridedot/crawler/engine"
	"github.com/stridedot/crawler/log"
	"github.com/stridedot/crawler/parse/douban"
	"github.com/stridedot/crawler/proxy"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
)

func main() {
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	proxyURLs := []string{"http://127.0.0.1:7890"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed", zap.Error(err))
	}

	cookie := "bid=zkq031h8Npg; __utmc=30149280; viewed=\"1007305_2089943\"; _pk_id.100001.8cb4=d560630ac0a50e46.1697943282.; douban-fav-remind=1; dbcl2=\"232581024:Ghreifrb9GA\"; ck=MoqN; frodotk_db=\"bd99ee018ffe63bf38cdca04affb2b3d\"; __utmz=30149280.1697956001.5.2.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; push_noty_num=0; push_doumail_num=0; __utmv=30149280.23258; _pk_ses.100001.8cb4=1; ap_v=0,6.0; __utma=30149280.1594306893.1697347189.1698413337.1698444472.13; __utmt=1; __utmb=30149280.38.2.1698444542323"
	var seeds []*collect.Request
	for i := 0; i <= 0; i += 25 {
		url := "https://www.douban.com/group/szsh/discussion?start=" + strconv.Itoa(i) + "&type=essence"
		seeds = append(seeds, &collect.Request{
			Url:       url,
			Cookie:    cookie,
			WaitTime:  2 * time.Second,
			ParseFunc: douban.ParseURL,
		})
	}

	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Logger:  logger,
		Proxy:   p,
	}

	s := engine.ScheduleEngine{
		WorkCount: 5,
		Logger:    logger,
		Fetcher:   &f,
		Seeds:     seeds,
	}
	s.Run()
}
