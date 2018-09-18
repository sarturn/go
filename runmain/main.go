package main

import (
	"fmt"
	"os"

	"github.com/woshimanong1990/go/spider"
)

func main() {
	url := "https://www.booktxt.net/2_3515/"
	dirPath := "./output"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
			return
		}
		fmt.Printf("mkdir success!\n")
	}

	urls, err := spider.Publisher(url)
	if err != nil {
		fmt.Println("spider start error", err)
		return
	}
	spider.Comsumer(urls, url)

}
