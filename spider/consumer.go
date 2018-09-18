package spider

import (
	"log"
	"os"
	"sync"

	"github.com/woshimanong1990/go/utils"
)

func WriteWithFileWrite(name, content string) {
	/*
		将解析的内容保存到文件中，如果文件不存在就创建
	*/
	contentByte, err := codetransutils.Decode([]byte(content)) //解码，不然中文乱码
	if err != nil {
		log.Fatal("Decode error", err)
		return
	}
	content = string(contentByte)
	fileObj, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("Failed to open the file", err.Error())
		return
		// os.Exit(2)
	}
	defer fileObj.Close()
	if _, err := fileObj.WriteString(content); err == nil {
		//fmt.Println("Successful writing to the file with os.OpenFile and *File.WriteString method.", content)
		return
	}
	// contents := []byte(content)
	// if _,err := fileObj.Write(contents);err == nil {
	//     fmt.Println("Successful writing to thr file with os.OpenFile and *File.Write method.",content)
	// }
}
func SpiderProcessor(cacheChan chan string, url string, waitGrop *sync.WaitGroup) error {
	/*
	  这里主要做爬取加解析html
	*/
	defer waitGrop.Done()
	// fmt.Println("url debug", url, name)
	content, err := Spider(url)
	if err != nil {
		log.Fatal(" consumer spider error:", url, err)
		return err
	}
	bodyContent, err := HtmlParseContent(content)
	if err != nil {
		log.Fatal(" consumer HtmlParseContent error:", url, err)
		return err
	}
	cacheChan <- bodyContent
	return nil
}
func WriterProcessor(cacheChan chan string, limitChan chan int, name string, waitGrop *sync.WaitGroup) error {
	/*
		这里主要是将解析的内容保存到文件
	*/
	defer waitGrop.Done()
	codeString, err := codetransutils.Decode([]byte(name))
	if err != nil {
		log.Fatal("decode title error", err)
		return err
	}
	name = "output/" + string(codeString) + ".txt"
	bodyContent := <-cacheChan
	WriteWithFileWrite(name, bodyContent)
	<-limitChan
	return nil
}
func Comsumer(urls map[string][]string, startUrl string) {
	/*
	   消费者，但是里面又有生产者消费者， 主要是加快爬取
	*/
	limitChan := make(chan int, 10)
	cacheChan := make(chan string, 10)
	var waitGrop sync.WaitGroup // 等待所有goroutine结束
	for _, value := range urls {
		url, name := value[0], value[1]
		url = startUrl + url
		limitChan <- 1
		waitGrop.Add(1)
		go SpiderProcessor(cacheChan, url, &waitGrop)
		waitGrop.Add(1)
		go WriterProcessor(cacheChan, limitChan, name, &waitGrop)

	}
	waitGrop.Wait()
}
