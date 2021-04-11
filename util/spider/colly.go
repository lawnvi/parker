package spider

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"log"
)

type HttpSpider struct {
	proxyIP []string
	Colly   *colly.Collector
}

func (c *HttpSpider) InitSpider() {
	c.Colly = colly.NewCollector()
	c.initRequest()
}

//init request
func (c *HttpSpider) initRequest() {
	c.Colly.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36")
	})
}

//设置代理ip池
//第三方服务 写一个ip池吧
//通过http获取
func (c *HttpSpider) SetProxyIps(ips []string) {
	c.proxyIP = ips
	rp, err := proxy.RoundRobinProxySwitcher(c.proxyIP...)
	if err != nil {
		log.Fatal(err)
	}
	c.Colly.SetProxyFunc(rp)
}
