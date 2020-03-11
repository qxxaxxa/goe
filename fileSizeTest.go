package main

import (

	"io/ioutil"
	"os"

)





func main() {
	s := "D:/Downloads/1.html"
	file, err := os.OpenFile(s, os.O_RDONLY, 6)
	if err != nil {
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()
	bytes, _ := ioutil.ReadAll(file)


}


