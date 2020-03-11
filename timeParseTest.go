package main

import (
	"fmt"
	"time"
)

func main() {

	loc, err := time.LoadLocation("UTC")
	fmt.Println(err)
	locTime, err := time.ParseInLocation("2006-01-02 15:04", "2020-02-21 10:21", loc)
	fmt.Println(err)
	fmt.Println(locTime)

}
