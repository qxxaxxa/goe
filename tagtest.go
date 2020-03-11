package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	file, _ := os.OpenFile("/Users/xiexingan/Downloads/tag1.html", os.O_RDWR, 6)
	bytes, _ := ioutil.ReadAll(file)
	reg, _ := regexp.Compile(`href="https://e-hentai.org/tag/([\w:+]+)`)
	allStringSubmatch := reg.FindAllStringSubmatch(string(bytes), -1)
	for _, submatch := range allStringSubmatch {
		fmt.Println(submatch[1])
	}

	split := strings.Split("kanimiso", ":")
	fmt.Println(split[0])
}
