package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var gidreg, _ = regexp.Compile(`var gid = (\d+);`)
var covarreg, _ = regexp.Compile(`background:transparent url\((.+)\) 0 0 no-repeat`)
var categoryReg, _ = regexp.Compile(`<div class="cs ct[123456789a]+" onclick="document.location='https://(?:exhentai|e-hentai).org/\w+'">(\w+)</div>`)
var titleReg, _ = regexp.Compile(`<h1 id="gn">(.+?)</h1>`)
var jtileReg, _ = regexp.Compile(`<h1 id="gj">(.+)</h1>`)
var tagReg, _ = regexp.Compile(`href="https://(?:e-hentai|exhentai).org/tag/([\w:+]+)`)
var uploaderReg, _ = regexp.Compile(`href="https://(?:e-hentai|exhentai)\.org/uploader/.+?">(.+?)(?:</a>)`)
var postedTimeReg, _ = regexp.Compile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}`)
var parentReg, _ = regexp.Compile(`<td class="gdt1">Parent:</td><td class="gdt2">(.+)</td>`)
var langReg, _ = regexp.Compile(`<td class="gdt1">Language:</td><td class="gdt2">(\w+) &nbsp;`)
var filesizeReg, _ = regexp.Compile(`<td class="gdt1">File Size:</td><td class="gdt2">([\d.]+ \w{2})</td>`)
var lenReg, _ = regexp.Compile(`<td class="gdt1">Length:</td><td class="gdt2">(\d+) pages</td>`)
var favReg, _ = regexp.Compile(`<td class="gdt2" id="favcount">(\d+) times</td>`)
var avgReg, _ = regexp.Compile(`<td id="rating_label" colspan="3">Average: ([\d.]+)</td>`)
var rtgCountReg, _ = regexp.Compile(`<span id="rating_count">(\d+)</span>`)
var chLoGrpReg, _ = regexp.Compile(`\[([^\[]+[汉漢]+[^]]*)]`)
var downurlReg, _ = regexp.Compile(`<a href="#" onclick="return popUp\('(.+)',480,320\)">Archive Download</a>`)

type galleryInfo struct {
	gid         int       `json:"gid"`
	url         string    `json:"url"`
	cover       string    `json:"cover"`
	category    string    `json:"category"`
	title       string    `json:"title"`
	jTitle      string    `json:"jTitle"`
	tags        []tag     `json:"tags"`
	uploader    string    `json:"uploader"`
	postedTime  time.Time `json:"postedTime"`
	parentId    int       `json:"parentId"`
	fileSize    float64   `json:"fileSize"`
	lang        string    `json:"lang"`
	lenth       int       `json:"lenth"`
	fav         int       `json:"fav"`
	avg         float64   `json:"avg"`
	rtgCount    int       `json:"rtgCount"`
	trGroup     string    `json:"trGroup"`
	downloadUrl string    `json:"downloadUrl"`
	chLoGrp     string    `json:"chLoGrp"`
}
type tag struct {
	tagType string
	title   string
}

const (
	b int64 = 1 << (10 * iota)
	kb
	mb
	gb
	tb
	pb
)
const (
	B  = "B"
	KB = "KB"
	MB = "MB"
	GB = "GB"
	TB = "TB"
	PB = "PB"
	EB = "EB"
)

var sizeMap = map[int]string{
	0: B, 1: KB, 2: MB, 3: GB, 4: TB, 5: PB, 6: EB,
}



func gtext() {
	//s := "/Users/xiexingan/anime/%d.html"
	s := "D:/Downloads/%d.html"
	infos := make([]galleryInfo, 0)
	for i := 0; i < 8; i++ {
		sprintf := fmt.Sprintf(s, i)
		file, _ := os.OpenFile(sprintf, os.O_RDONLY, 6)
		defer func() {
			if err := file.Close(); err != nil {
				return
			}
		}()
		bytes, _ := ioutil.ReadAll(file)
		text := string(bytes)

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
			url:         sprintf,
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
		infos = append(infos, gi)

	}

	sum := 0.0
	for _, info := range infos {
		sum+=info.fileSize
	}
	fmt.Println(fitAmout(sum))

}

func (gi galleryInfo) String() string {

	fita, unit := fitAmout(gi.fileSize)
	a := "汉化组:" + gi.chLoGrp + "\n" +
		"下载链接:" + gi.downloadUrl + "\n" +
		"平均分:" + strconv.FormatFloat(gi.avg, 'f', -1, 64) + "\n" +
		"分类:" + gi.category + "\n" +
		"封面:" + gi.cover + "\n" +
		"收藏数:" + strconv.Itoa(gi.fav) + "\n" +
		"大小:" + strconv.FormatFloat(fita, 'f', -1, 64) + unit + "\n" +
		"id:" + strconv.Itoa(gi.gid) + "\n" +
		"日文标题:" + gi.jTitle + "\n" +
		"语言:" + gi.lang + "\n" +
		"页数:" + strconv.Itoa(gi.lenth) + "\n" +
		"父id:" + strconv.Itoa(gi.parentId) + "\n" +
		"上传时间:" + gi.postedTime.String() + "\n" +
		"评价人数:" + strconv.Itoa(gi.rtgCount) + "\n" +
		"tags:" + strconv.Itoa(len(gi.tags)) + "\n" +
		"标题:" + gi.title + "\n" +
		"汉化组:" + gi.trGroup + "\n" +
		"up:" + gi.uploader + "\n" +
		"url:" + gi.url + "\n"
	return a
}

func getPostTime(text string) time.Time {
	allString := postedTimeReg.FindAllString(text, -1)
	timeText := allString[0]
	loc, _ := time.LoadLocation("UTC")
	pTime, _ := time.ParseInLocation("2006-01-02 15:04", timeText, loc)
	return pTime
}

func getFileSize(text string) float64 {
	var f float64
	filesizetext := getRegMatch(text, filesizeReg)
	split := strings.Split(filesizetext, " ")
	size, _ := strconv.ParseFloat(split[0], 64)
	switch split[1] {
	case B:
		f = size * float64(b)
	case KB:
		f = size * float64(kb)
	case MB:
		f = size * float64(mb)
	case GB:
		f = size * float64(gb)
	case TB:
		f = size * float64(tb)
	}
	return f
}

func tagsBuild(tags []tag, text string) []tag {
	tags = make([]tag, 256)
	tagMatch := tagReg.FindAllStringSubmatch(text, -1)
	for i, submatch := range tagMatch {
		full := submatch[1]
		var t tag
		split := strings.Split(full, ":")
		if strings.Contains(full, ":") {
			t = tag{tagType: split[0], title: split[1],}
		} else {
			t = tag{tagType: "misc", title: split[0],}
		}
		tags[i] = t
	}
	return tags
}

func getRegMatch(text string, reg *regexp.Regexp) (re string) {
	submatch := reg.FindAllStringSubmatch(text, -1)
	for _, ss := range submatch {
		re = ss[1]
		break
	}
	return re
}

func fitAmout(a float64) (float64, string) {
	for i := 0; i < 10; i++ {
		s := 1 << (10 * i)

		if int(a)/s < 1024 {
			re := a / float64(s)
			return re, sizeMap[i]
		}

	}
	return 0, ""
}
