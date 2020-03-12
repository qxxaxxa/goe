package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func checksha1() {
	base := cachePath()
	fmt.Println("start scan saved image file")

	files, _ := ioutil.ReadDir(base)
	for _, dir := range files {
		if dir.IsDir() {
			fmt.Println(dir.Name())
			path := filepath.Join(base, dir.Name())
			fileInfos, _ := ioutil.ReadDir(path)
			for _, info := range fileInfos {
				if info.IsDir() {
					images, _ := ioutil.ReadDir(filepath.Join(base, dir.Name(), info.Name()))
					for _, image := range images {
						time.Sleep(time.Millisecond*1)
						go readAndSum(base, dir, info, image)
					}

				}
			}
		}
	}
	fmt.Println(" scan file complete")
	fmt.Println("start read  unsaved image list")
}

func readAndSum(base string, dir os.FileInfo, info os.FileInfo, image os.FileInfo) {
	st := time.Now()
	imagePath := filepath.Join(base, dir.Name(), info.Name(), image.Name())
	file, _ := os.OpenFile(imagePath, os.O_RDONLY, os.ModePerm)
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()
	bytes, _ := ioutil.ReadAll(file)
	hash := sha1.New()
	hash.Write(bytes)
	sum := hash.Sum(nil)
	en := time.Now()
	fmt.Println(image.Name(), "comsume:", en.Sub(st).Milliseconds(), "millsencond")
	if fmt.Sprintf("%x", sum) != image.Name()[0:40] {
		fmt.Println(imagePath)
	}
}
