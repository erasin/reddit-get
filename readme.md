# reddit-get

> 查看 <https://erasin.wang/go-colly/>

使用 golang 来创建一个爬虫获取 <[reddit.com](http://www.reddit.com)> 图片。

比如 [r/wallpaper](http://www.reddit.com/r/wallpaper),通过解析官方API <http://www.reddit.com/r/wallpaper.json?limit=22&after=xxxxx> 返回的 `JSON` 数据来分析和下载文件。

主要使用库

- **github.com/urfave/cli/v2** 用来创建命令行。
- **github.com/gocolly/colly**  是golang爬虫框架, 用来获取数据。
- **github.com/buger/jsonparser**  来解析 reddit 的json数据。
