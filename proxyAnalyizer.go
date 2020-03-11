package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)


func ana() {
	reg, _ := regexp.Compile(`([\d]{4}-[\d]{2}-[\d]{2}T[\d]{2}:[\d]{2}):[\d]{2}Z \[debug] Proxy file download request complete for [\w]{40}-(\d+)-`)
	proreg, _ := regexp.Compile(`([\d]{4}-[\d]{2}-[\d]{2}T[\d]{2}:[\d]{2}):[\d]{2}Z \[info] {[\d]+\/[\d.]+}[\s]+Code=200 Bytes=([\d]+)[\s]+Finished processing request in `)
	file, _ := os.OpenFile(logPath(), os.O_RDWR, 6)

	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	bytes, _ := ioutil.ReadAll(file)
	var smap = make(map[string]int)
	allString := reg.FindAllStringSubmatch(string(bytes), -1)
	for _, ss := range allString {
		atoi, _ := strconv.Atoi(ss[2])
		smap[ss[1]] += atoi
	}
	sum := 0
	for _, i := range smap {
		sum += i
		//fmt.Println(s, "cache in:", i/1024/60)

	}
	fmt.Println("cache in", sum/len(smap)/1024/1024, "MB/min")
	pmap := make(map[string]int)
	cmap := make(map[string]int)
	submatch := proreg.FindAllStringSubmatch(string(bytes), -1)
	for _, strings := range submatch {
		atoi, _ := strconv.Atoi(strings[2])
		cmap[strings[1]] += 1
		pmap[strings[1]] += atoi
	}
	for s, i := range pmap {
		fmt.Println(s, "process in:", i/1024/60, "KB/s")
	}
	count := 0
	for s, i := range cmap {
		count += i
		fmt.Println(s, "process in:", i, "Times")
	}
	fmt.Println("avg count:", count/len(cmap), "times/ min")
}