package main

import (
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

var pageUrlChan = make(chan string, 1024)
var pageDataChan = make(chan []byte, 1024)
var cli *http.Client
var glchan = make(chan string, 1024)

const pageFormat = "https://exhentai.org/?page=%d"

/*
set http_proxy=http://127.0.0.1:7890
set https_proxy=http://127.0.0.1:7890

*/

func crawler() {
	proxyAddr := "http://127.0.0.1:7890"
	cli = newHttpProxyCli(proxyAddr)
	//lastPage, done := getPageCount()
	//if done {
	//	return
	//}
	quit := make(chan bool)
	os.Create(galleryPath())
	lastPage := 10
	go func() {
		for i := 1; i < lastPage+1; i++ {
			go formatPageUrl(i)
		}
	}()

	go func() {
		for pageUrl := range pageUrlChan {
			go func() {
				data := getPageData(pageUrl, func() {
					pageUrlChan <- pageUrl
				})
				if data != nil {
					pageDataChan <- data
				}
			}()

		}
	}()

	gregex, _ := regexp.Compile(`https://exhentai\.org/g/\d+/\w+/`)
	fmt.Println("compile success")
	go func() {
		for bytes := range pageDataChan {
			go extractGalleryUrl(gregex, bytes)
		}

	}()
	go func() {
		fmt.Println("begin open file")

		file, err := os.OpenFile(galleryPath(), os.O_RDWR, 6)

		defer func() {
			if err := file.Close(); err != nil {
				fmt.Println("file close err")
				return
			}
		}()
		if err != nil {
			fmt.Println("open file err")
			return
		}

		for gurl := range glchan {
			writeGalleryUrltoFile(file, gurl)
			data := getPageData(gurl, func() {
				glchan <- gurl
			})
			text := string(data)
			var (
				gid         int
				trGroup     string
				downloadUrl string
				cover       string
				category    string
				title       string
				jTitle      string
				uploader    string
				pTime       time.Time
				tags        []tag
				parentId    int
				fileSize    float64
				lang        string
				lenth       int
				fav         int
				avg         float64
				rtgCount    int
				chLoGrp     string
			)
			gid, _ = strconv.Atoi(getRegMatch(text, gidreg))
			parentId, _ = strconv.Atoi(getRegMatch(text, parentReg))
			cover = getRegMatch(text, covarreg)
			category = getRegMatch(text, categoryReg)
			title = getRegMatch(text, titleReg)
			jTitle = getRegMatch(text, jtileReg)
			tags = tagsBuild(tags, text)
			uploader = getRegMatch(text, uploaderReg)
			pTime = getPostTime(text)

			fileSize = getFileSize(text)
			lang = getRegMatch(text, langReg)
			lenth, _ = strconv.Atoi(getRegMatch(text, lenReg))
			fav, _ = strconv.Atoi(getRegMatch(text, favReg))
			avg, _ = strconv.ParseFloat(getRegMatch(text, avgReg), 64)
			rtgCount, _ = strconv.Atoi(getRegMatch(text, rtgCountReg))
			downloadUrl = html.UnescapeString(getRegMatch(text, downurlReg))

			trGroup = getRegMatch(text, chLoGrpReg)
			gi := galleryInfo{
				gid:         gid,
				url:         gurl,
				cover:       cover,
				category:    category,
				title:       title,
				jTitle:      jTitle,
				tags:        tags,
				uploader:    uploader,
				postedTime:  pTime,
				parentId:    parentId,
				fileSize:    fileSize,
				lang:        lang,
				lenth:       lenth,
				fav:         fav,
				avg:         avg,
				rtgCount:    rtgCount,
				trGroup:     trGroup,
				downloadUrl: downloadUrl,
				chLoGrp:     chLoGrp,
			}
			fmt.Println(gi)

		}
	}()

	<-quit

}
func writeGalleryUrltoFile(file *os.File, gurl string) {
	offset, _ := file.Seek(0, io.SeekEnd)
	_, err := file.WriteAt([]byte(gurl+"\n"), offset)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func getPageData(pageUrl string, callback func()) []byte {
	start := time.Now()

	bytes, err := proxyHttpGet(pageUrl)
	end := time.Now()
	fmt.Println("花费时间:", end.Sub(start).Milliseconds())
	if err != nil {
		fmt.Println("失败url:", pageUrl)
		go callback()
		return nil
	}
	return bytes

}

func formatPageUrl(i int) {
	pageUrlChan <- fmt.Sprintf(pageFormat, i)
}

func savetofile() {
	fmt.Println("begin open file")
	file, err := os.OpenFile(galleryPath(), os.O_RDWR, 6)

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("file close err")
			return
		}
	}()
	if err != nil {
		fmt.Println("open file err")
		return
	}

	for gurl := range glchan {
		offset, _ := file.Seek(0, io.SeekEnd)
		_, err := file.WriteAt([]byte(gurl+"\n"), offset)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("写入成功")

	}

}

func extractGalleryUrl(gregex *regexp.Regexp, bytes []byte) {
	urls := gregex.FindAllString(string(bytes), -1)
	for _, s := range urls {
		glchan <- s
	}
}

func getPageCount() (int, bool) {
	ehentaiIndexpageUrl := "https://exhentai.org/"
	indexDate, err := proxyHttpGet(ehentaiIndexpageUrl)
	if err != nil {
		fmt.Println("get index failed")
		return 0, true
	}
	lastPage := getlastPage(string(indexDate))
	return lastPage, false
}

func getlastPage(s string) int {

	reg, _ := regexp.Compile(`Jump to page: \(1-(\d+)\)`)
	submatch := reg.FindAllStringSubmatch(s, 1)
	i, err := strconv.Atoi(submatch[0][1])
	if err != nil {
		fmt.Println("read page failed")
		return 0
	}
	return i

}

func newHttpProxyCli(proxyAddr string) *http.Client {
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		fmt.Println("url pas err")
		return nil
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
		DialContext: (&net.Dialer{
			Timeout:   time.Second * 15,
			KeepAlive: time.Second * 30,
		}).DialContext,
		TLSHandshakeTimeout:   time.Second * 10,
		MaxIdleConns:          100,
		IdleConnTimeout:       time.Second * 100,
		ExpectContinueTimeout: time.Second * 1,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   time.Second * 15,
	}

}

func proxyHttpGet(url string) ([]byte, error) {
	req, err := buildRequest(url)
	if err != nil {
		return nil, err
	}
	response, err := cli.Do(req)
	defer func() {
		if response != nil {
			if qerr := response.Body.Close(); qerr != nil {
				fmt.Println("close body err")
				return
			}
		}

	}()
	if err != nil {
		fmt.Println("do err:", err)
		return nil, err
	}

	bytes, err := ioutil.ReadAll(response.Body)

	fmt.Printf("read %d bytes\n", len(bytes))
	if err != nil {
		fmt.Println("read err")
		return nil, err
	}
	return bytes, nil
}

func buildRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("new req err")
		return nil, err
	}

	domain := ".exhentai.org"
	req.AddCookie(&http.Cookie{
		Name:   "hath_perks",
		Value:  "m1.a.t1.p1.p2.p3.s-0d165fd7c5",
		Path:   "/",
		Domain: domain,
	})
	req.AddCookie(&http.Cookie{
		Name:   "igneous",
		Value:  "3291eecbf",
		Path:   "/",
		Domain: domain,
		Secure: false,
	})
	req.AddCookie(&http.Cookie{
		Name:   "ipb_member_id",
		Value:  "1619825",
		Path:   "/",
		Domain: domain,
	})
	req.AddCookie(&http.Cookie{
		Name:   "ipb_pass_hash",
		Value:  "a90e0df06e41719089b1c06743e0539a",
		Path:   "/",
		Domain: domain,
	})
	req.AddCookie(&http.Cookie{
		Name:   "sk",
		Value:  "3t1fry8ovjo1nxlcjmdh0jn6qrdk",
		Path:   "/",
		Domain: domain,
	})
	req.AddCookie(&http.Cookie{
		Name:   "yay",
		Value:  "louder",
		Path:   "/",
		Domain: domain,
	})
	return req, nil
}
