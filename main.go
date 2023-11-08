package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"

	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// 设置代理
var (
	cookie = ""
	// 失败次数
	failedNum = 0

	sw            = sync.WaitGroup{}
	dirLoc        = "./pixiv/"
	proxyUrl,_ = url.Parse("http://127.0.0.1:7890")

	client = http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
)

func headerSet(req *http.Request) {
	req.Header.Set("authority", "www.pixiv.net")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "zh,zh-CN;q=0.9")
	req.Header.Set("cache-control", "max-age=0")
	
	req.Header.Set("cookie", cookie)
	
	req.Header.Set("referer", "https://accounts.pixiv.net/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="116", "Not)A;Brand";v="24", "Google Chrome";v="116"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
		// 增加请求来源
	req.Header.Add("referer", "https://www.pixiv.net/")
}

// 得到图片pid
func solveRankPage(rankPage string) string {

	// 进入排行榜页面，发送请求
	req, err := http.NewRequest("GET", rankPage, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 设置header信息
	headerSet(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("resp.Body read error ", err)
	}
	// 解析页面，获取pid
	reg, err := regexp.Compile("(/artworks/\\d+)")
	if err != nil {
		fmt.Println("regexp compile error ", err)
	}
	pid := reg.FindString(string(bodyText))
	fmt.Printf("pid is %v\n", pid)

	return "https://www.pixiv.net" + pid
}

// DownloadFile 下载文件
func DownloadFile(url, address, fileName string) (ok bool) {
	fileName = address + fileName
	_,err := os.Stat(fileName)

	if err == nil {
		fmt.Println("该文件已存在")
		return true;
	}


	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("NewRequest ", err)
	}
	headerSet(req)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	

	err = os.WriteFile(fileName, bytes, 0666)
	if err != nil {
		fmt.Println("文件保存错误 ", err)
		return false
	} else {
		return true
	}
}

// 获取图片链接
func solveImgLink(link string) (string, string, error) {

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 设置header信息
	headerSet(req)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("resp.body read error ", err)
		return "nil","nil",err
	}
	var ImgLink string

	reg, err := regexp.Compile("https://i.pximg.net/img-original/img/\\d{4}/\\d{2}/\\d{2}/\\d{2}/\\d{2}/\\d{2}/\\w+\\.(jpg|jpeg|png|gif)")
	if err != nil {
		fmt.Println("regexp compile error ", err)

		return "nil","nil",err
	}
	ImgLink = reg.FindString(string(bodyText))

	fmt.Println("ImgLink is ", ImgLink)

	urlSp := strings.Split(ImgLink, "/")
	if len(urlSp) < 11 {
		fmt.Println("urlSp error ", ImgLink)
		
		return "nil","nil",errors.New("url error")
	}
	fileName := fmt.Sprintf("%s_%s_%s_%s_%s_%s_pid_%s", urlSp[5], urlSp[6], urlSp[7], urlSp[8], urlSp[9], urlSp[10], urlSp[11])

	return ImgLink, fileName, err
}

func task(start int, work int, days int, target string) {
	failed := 0
	for i := start; i < start+work && i <= days; i++ {
		time.Sleep(time.Millisecond * (1000 * 2) )
		visitTime := GetDate(i)
		fmt.Println(visitTime)
		link := solveRankPage(target + visitTime)
		fmt.Println("link is ", link)
		// 进入详细页面
		imgLink, filename,err := solveImgLink(link)
		if err != nil {
			fmt.Println("solveImgLink error",err);
			failed++
			continue
		}
		path := Mkdir(visitTime, dirLoc)
		if DownloadFile(imgLink, path, filename) == true {
			fmt.Println("下载成功！文件名为", filename)
			fmt.Println("下载目录为", path)
		}
	}
	fmt.Println("----- Task Over -----")
	sw.Done()
	failedNum += failed
}

func rangeTime(days int, target string) {

	days += 1	
	p := days / 5
	for i := 2; i <= days ; i += p {
		sw.Add(1)
		go task(i, p, days, target)
	}
	sw.Wait()

}



func readCookie(cookiePath string) (string,error) {
	if cookiePath == "" {
		return "",errors.New("empty args")
	}
	file, err := os.Open(cookiePath)
	if err != nil {
		fmt.Println("无法打开文件",err)
		return "",errors.New("file error")
	}
	defer file.Close()

	// 使用bufio包创建一个Scanner来逐行读取文件内容
	scanner := bufio.NewScanner(file)
	cookie := ""
	for scanner.Scan() {
		cookie += scanner.Text()
	}


	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件时发生错误:", err)
		return "",err
	}
	return cookie,nil

}


func main() {

	var days int
	var cookiePath string

	flag.IntVar(&days, "d", 100, "爬取天数")
	flag.StringVar(&cookiePath,"c","","cookie")
	// 解析命令行参数
	flag.Parse()

	c,err := readCookie(cookiePath)	
	if (err != nil) {
		fmt.Println("cookie get failed! ",err)
		return
	}
	cookie = c
	
	// 创建pixiv文件夹
	if pathExist(dirLoc) == false {
		err := os.Mkdir(dirLoc,0755) 
		if err != nil {
			fmt.Println("文件夹创建失败",err)
		}
	}

	start := time.Now()
	rankPage := "https://www.pixiv.net/ranking.php?mode=daily&date=" // + now
	rangeTime(days, rankPage)
	elapsed := time.Since(start)
	fmt.Println("程序耗时: ", elapsed)
	fmt.Println("失败次数：", failedNum);
}
