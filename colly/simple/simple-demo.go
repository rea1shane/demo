package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
)

func main() {

	/**
	首先，调用 colly.NewCollector() 创建一个类型为 *colly.Collector 的爬虫对象。
	由于每个网页都有很多指向其他网页的链接。如果不加限制的话，运行可能永远不会停止。
	所以上面通过传入一个选项 colly.AllowedDomains("www.baidu.com") 限制只爬取域名为 www.baidu.com 的网页。
	*/
	c := colly.NewCollector(
		colly.AllowedDomains("www.baidu.com"),
	)

	/**
	然后我们调用 c.OnHTML 方法注册 HTML 回调，对每个有 href 属性的 a 元素执行回调函数。
	这里继续访问 href 指向的 URL。也就是说解析爬取到的网页，然后继续访问网页中指向其他页面的链接。
	*/
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	/**
	调用 c.OnRequest() 方法注册请求回调，每次发送请求时执行该回调，这里只是简单打印请求的 URL。
	*/
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	/**
	调用 c.OnResponse() 方法注册响应回调，每次收到响应时执行该回调，这里也只是简单的打印 URL 和响应大小。
	*/
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Response %s: %d bytes\n", r.Request.URL, len(r.Body))
	})

	/**
	调用 c.OnError() 方法注册错误回调，执行请求发生错误时执行该回调，这里简单打印 URL 和错误信息。
	*/
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error %s: %v\n", r.Request.URL, err)
	})

	/**
	最后我们调用 c.Visit() 开始访问第一个页面。
	*/
	c.Visit("http://www.baidu.com/")

}
