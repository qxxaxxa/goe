package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {

	base := cachePath()
	files, _ := ioutil.ReadDir(base)
	imap := make(map[string]int)
	for _, dir := range files {
		if dir.IsDir() {
			path := filepath.Join(base, dir.Name())
			fileInfos, _ := ioutil.ReadDir(path)
			for _, info := range fileInfos {
				if info.IsDir() {
					images, _ := ioutil.ReadDir(filepath.Join(base, dir.Name(), info.Name()))
					for _, image := range images {
						imap[strings.Split(image.Name(), "-")[0]] = 1
					}

				}
			}
		}
	}

	file, _ := os.OpenFile(imagePath(), os.O_RDONLY, 6)
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()
	all, _ := ioutil.ReadAll(file)
	reg, _ := regexp.Compile(`\./[\dabcde]{2}/[\dabcde]{2}/([\w]+-[\d]+-[\d]+-[\d]+-[\w]+)`)
	submatch := reg.FindAllStringSubmatch(string(all), -1)
	for i, ss := range submatch {
		fileid := ss[1]
		split := strings.Split(fileid, "-")
		if imap[split[0]] == 1 {
			continue
		}
		fmt.Printf("第%d个\n", i)
		rand.Seed(time.Now().UnixNano())
		fileindex := strconv.Itoa(rand.Intn(50000))
		keystampTime := strconv.FormatInt(time.Now().Unix(), 10)
		clientKey := "X72CuATxVA6Yhy4sTmws"
		sum := sha1.Sum([]byte(keystampTime + "-" + fileid + "-" + clientKey + "-hotlinkthis"))
		urlstring := "https://" + hostprefix() + ".ehedgzdwvjcc.hath.network:11759/h/" + fileid + "/keystamp=" + keystampTime + "-" + fmt.Sprintf("%x", sum)[0:10] + ";fileindex=" + fileindex + ";xres=org/" + split[0] + "." + split[4]
		fmt.Println(urlstring)

		download(urlstring, base, split, fileid)
	}

	quit := make(chan int)
	<-quit

}

func download(urlstring string, base string, split []string, fileid string) {
	func() {
		response, _ := http.Get(urlstring)
		defer func() {
			if err := response.Body.Close(); err != nil {
				return
			}
		}()

		bytes, _ := ioutil.ReadAll(response.Body)

		path01 := filepath.Join(base, split[0][0:2])
		notExistorCreat(path01)
		path02 := filepath.Join(path01, split[0][2:4])
		notExistorCreat(path02)
		//join := filepath.Join(base, fileid)
		join := filepath.Join(path02, fileid)
		tfile, err := os.Create(join)
		if err != nil {
			fmt.Println(err)
			return
		}
		n, _ := tfile.Write(bytes)
		fmt.Println("写入", n)
	}()
}

func notExistorCreat(path01 string) bool {
	_, err2 := os.Stat(path01)
	if os.IsNotExist(err2) {
		if err := os.Mkdir(path01, os.ModePerm); err != nil {
			return true
		}
	}
	return false
}
