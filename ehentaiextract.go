package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

func extract() {
	reg, _ := regexp.Compile("Code=200.+(/h.+) HTTP/1.1")
	file, _ := os.OpenFile(logPath(), os.O_RDWR, 6)
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()
	bytes, _ := ioutil.ReadAll(file)
	allString := reg.FindAllStringSubmatch(string(bytes), -1)

	hurl := "https://" + hostprefix() + "." + hentaihost() +
		".hath.network:11759"

	linkfile, _ := os.Create(linkPath())
	defer func() {
		if err := linkfile.Close(); err != nil {
			return
		}
	}()
	reg2, _ := regexp.Compile("/h/(.{40})-")
	mmap := make(map[string]string)
	for _, strings := range allString {
		submatch := reg2.FindAllStringSubmatch(strings[1], -1)

		for _, i := range submatch {
			mmap[i[1]] = strings[1]

		}

	}
	for _, s := range mmap {
		wurl := hurl + s + "\n"
		ret, _ := linkfile.Seek(0, io.SeekEnd)
		linkfile.WriteAt([]byte(wurl), ret)
		fmt.Println(wurl)
	}

	//for _, strings := range allString {
	//	wurl := hurl + strings[1]+"\n"
	//	ret, _ := linkfile.Seek(0, io.SeekEnd)
	//	linkfile.WriteAt([]byte(wurl), ret)
	//	fmt.Println(wurl)
	//}
}
