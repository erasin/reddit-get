package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

// DownloadFile 需要下载的文件
type DownloadFile struct {
	Filename string
	Folder   string
	URL      string
}

// Down 开始下载
func (f *DownloadFile) Down(directory string) {
	wg.Add(1)
	defer wg.Done()

	// 实际路径
	p := path.Join(directory, f.Folder, f.Filename)

	if _, err := os.Stat(p); err == nil {
		log.Println("文件已存在:", f.URL)
		return
	}
	os.Mkdir(path.Join(directory, f.Folder), 0777)

	output, err := os.Create(p)
	defer output.Close()
	if err != nil {
		log.Println("创建失败: ", err)
	}

	response, err := http.Get(f.URL)
	if err != nil {
		log.Println("下载失败: ", err)
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		log.Println("写入失败 ", err)
	}
	log.Printf("下载文件 %s/%s", f.Folder, f.Filename)

}
