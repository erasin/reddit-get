package main

import (
	"fmt"
	"html"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/buger/jsonparser"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"github.com/urfave/cli/v2"
)

var nowPage int64   // 当前页
var limitPage int64 // 限制数量
var startPage int64 // 限制数量

func main() {
	app := &cli.App{
		Name:  "reddit-get",
		Usage: "download image!",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "chanel",
				Aliases: []string{"c"},
				Value:   "wallpaper",
				Usage:   "频道名称",
			},
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Value:   "output",
				Usage:   "输出目录",
			},
			&cli.Int64Flag{
				Name:        "start",
				Aliases:     []string{"s"},
				Value:       1,
				Usage:       "start page",
				Destination: &startPage,
			},
			&cli.Int64Flag{
				Name:        "limit",
				Aliases:     []string{"l"},
				Value:       10,
				Usage:       "limit page",
				Destination: &limitPage,
			},
		},
		Action: func(c *cli.Context) error {
			run(c.String("chanel"), c.String("dir"))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var wg sync.WaitGroup // 计数器
var limit int64 = 24  // reddit 每页数量

func run(chanel, dir string) {
	log.Printf("--> %s/%s\n", dir, chanel)
	endPage := startPage + limitPage
	// 创建默认
	c := colly.NewCollector(colly.AllowURLRevisit())

	// 加载代理
	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:7891")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

	// 记录
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(nowPage, " --> 访问：", r.URL.String())
	})

	// 检查元素
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// 相应处理
	c.OnResponse(func(r *colly.Response) {
		nowPage++

		// 后缀
		fileExtToDownload := map[string]bool{".jpg": true, ".png": true, ".gif": true}

		if nowPage >= startPage {
			// 解析
			jsonparser.ArrayEach(r.Body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				url, _, _, _ := jsonparser.Get(value, "data", "url")
				formatedURL := html.UnescapeString(string(url))
				_, filename := filepath.Split(formatedURL)
				if fileExtToDownload[filepath.Ext(filename)] {
					f := DownloadFile{Filename: filename, Folder: chanel, URL: formatedURL}
					go f.Down(dir)
				}
			}, "data", "children")

			// 处理数量
			limit, err = jsonparser.GetInt(r.Body, "data", "dist")
			if err != nil {
				limit = 24
			}
		}

		// 继续处理
		before, err := jsonparser.GetString(r.Body, "data", "after")
		if err == nil && nowPage < endPage {
			u := fmt.Sprintf("http://www.reddit.com/r/%s.json?limit=%d&after=%s", chanel, limit, before)
			wg.Add(1)
			go c.Visit(u)
		}

		// 结束 visit
		wg.Done()
	})

	u := fmt.Sprintf("http://www.reddit.com/r/%s.json", chanel)
	wg.Add(1) // +1
	go c.Visit(u)

	wg.Wait() // 等待清零
}

// https://plantuml.com/es/code-php
