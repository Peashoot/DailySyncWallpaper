package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// 获取图片url
func getImgURL() (bool, string) {
	now := time.Now().Unix()
	url := fmt.Sprintf("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&nc=%d&pid=hp", now)
	Println("start to get " + url)
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	var imgURL string
	if err != nil {
		Fatalln(err)
		return false, imgURL
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	Println(fmt.Sprintf("The response's status code is %d", res.StatusCode))
	result := res.StatusCode == 200
	// 如果请求成功，解析json获取到imgUrl
	if result {
		var dat map[string]interface{}
		json.Unmarshal(body, &dat)
		if v, ok := dat["images"]; ok {
			ws := v.([]interface{})
			for _, wsItem := range ws {
				wsMap := wsItem.(map[string]interface{})
				if vCw, ok := wsMap["url"]; ok {
					imgURL = fmt.Sprintf("https://cn.bing.com%s", vCw)
					Println("Success get url of image -- " + imgURL)
					break
				}
			}
		}
	}
	return result, imgURL
}

// 从url中获取图片名称
func getFileName(imgURL string) string {
	u, err := url.Parse(imgURL)
	if err != nil {
		Fatalln(err)
	}
	m, _ := url.ParseQuery(u.RawQuery)
	return m["id"][0]
}

// 保存网络图片到本地
func saveImgs(imgURL string, fileName string) {
	res, err := http.Get(imgURL)
	if err != nil {
		Fatalln(err)
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	file, err := os.Create(fileName)
	if err != nil {
		Fatalln(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	io.Copy(writer, reader)
	// 释放图片资源
	defer file.Close()
}

// 下载bing图片并设置为壁纸
func download() {
	result, imgURL := getImgURL()
	if result {
		fileName := config.ImgPath + getFileName(imgURL)
		//如果图片文件不存在，就将图片下载下来
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			Println(fileName + " doesn't exist, start download...")
			saveImgs(imgURL, fileName)
			Println(fileName + " has been saved in current directory")
			Println("Start set wallpaper...")
			setWallpaper(fileName)
			Println("End set wallpaper...")
		} else {
			Println(fmt.Sprintf("Image %s already exists...", fileName))
		}
	}
}
