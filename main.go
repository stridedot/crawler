package main

import (
	"fmt"
	"github.com/stridedot/crawler/collect"
	"github.com/stridedot/crawler/log"
	"github.com/stridedot/crawler/parse/douban"
	"github.com/stridedot/crawler/proxy"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
	"time"
)

func main() {
	plugin, closer := log.NewFilePlugin("./test.log", zapcore.InfoLevel)
	defer closer.Close()
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	proxyURLs := []string{"http://127.0.0.1:7890"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed", zap.Error(err))
	}

	cookie := "bid=zkq031h8Npg; __utmc=30149280; viewed=\"1007305_2089943\"; _pk_id.100001.8cb4=d560630ac0a50e46.1697943282.; douban-fav-remind=1; _pk_ses.100001.8cb4=1; dbcl2=\"232581024:Ghreifrb9GA\"; ck=MoqN; frodotk_db=\"bd99ee018ffe63bf38cdca04affb2b3d\"; __utma=30149280.1594306893.1697347189.1697943283.1697956001.5; __utmz=30149280.1697956001.5.2.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmt_douban=1; push_noty_num=0; push_doumail_num=0; __utmt=1; __utmv=30149280.23258; __utmb=30149280.3.10.1697956001"
	var workList []*collect.Request
	for i := 0; i <= 0; i += 25 {
		url := "https://www.douban.com/group/szsh/discussion?start=" + strconv.Itoa(i)
		workList = append(workList, &collect.Request{
			Url:       url,
			Cookie:    cookie,
			ParseFunc: douban.ParseURL,
		})
	}

	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}

	for len(workList) > 0 {
		items := workList
		workList = nil
		for _, item := range items {
			body, err := f.Get(item)
			time.Sleep(3 * time.Second)
			if err != nil {
				logger.Error("read content failed",
					zap.Error(err),
				)
				continue
			}
			res := item.ParseFunc(body)
			for _, value := range res.Items {
				fmt.Printf("%s\n", value.(string))
				logger.Info("result",
					zap.String("item", value.(string)),
				)
			}
			workList = append(workList, res.Requests...)
		}
	}
}
