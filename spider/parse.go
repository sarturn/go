package spider

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func renderNode(n *html.Node) string {
	/*
		渲染html内容，或者说是转成字符串
	*/
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func ParseContent(doc *html.Node) (string, error) {
	/*
		解析小说内容，其实是做了一个替换，也许有更好的方法
	*/
	var content []string
	for i := doc.FirstChild; i != nil; i = i.NextSibling {
		if i.Data == "br" {
			content = append(content, "\r\n")
		} else {
			parseContent := renderNode(i)
			hexStr := "c2a0"
			replaceString, err := hex.DecodeString(hexStr)
			if err != nil {
				fmt.Println("hex decode error", err)
				return "", err
			}
			parseContent = strings.Replace(parseContent, string(replaceString), " ", 4) // 前面有几个空格，中文转码时会有问题，替换下
			content = append(content, parseContent)
		}
	}
	return strings.Join(content, ""), nil
}
func ParseTitle(n *html.Node) string {
	/*
		解析章节标题
	*/
	var title string = ""
	if n.Type == html.ElementNode && n.Data == "div" && title == "" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "bookname" {
				for i := n.FirstChild; i != nil; i = i.NextSibling {

					if i.Data == "h1" {
						title = "   " + renderNode(i.FirstChild)
						// fmt.Println("title:", title)
						break
					}

				}
				break
			}
		}
	}
	return title
}
func HtmlParseContent(content string) (string, error) {
	/*
		解析html内容
	*/
	var bodyContent string = ""
	var title string = ""
	var haveFound bool = false
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if title == "" {
			title = ParseTitle(n)
		}
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == "content" {
					bodyContent, _ = ParseContent(n)
					// fmt.Println("title:", title)
					bodyContent = title + "\r\n" + bodyContent
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if haveFound {
				break
			}
			f(c)
		}
	}
	f(doc)
	return bodyContent, nil
}
func HtmlParseChapter(content string) (map[string][]string, error) {
	/*
		解析各个章节，返回的是个map，key为章节序号，value是slice，第一个是url，第二个是第几章，第三个是章节标题
	*/
	var bodyContent string = ""
	var f func(*html.Node)
	var f2 func(*html.Node)
	var haveFoundDiv bool = false

	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	urls := make(map[string][]string)
	f2 = func(n *html.Node) {
		var haveFoundA bool = false
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					bodyContent = renderNode(n.FirstChild)
					arr := strings.SplitN(bodyContent, ".", 2)
					charArray := make([]string, 1)
					charArray[0] = a.Val
					charArray = append(charArray, strings.SplitN(arr[1], " ", 2)...)
					urls[arr[0]] = charArray
					// fmt.Println("debug", renderNode(n.FirstChild))
					// fmt.Println("debug", renderNode(n))
					haveFoundA = true
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if haveFoundA {
				break
			}
			f2(c)
		}
	}
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == "list" {
					f2(n)
					haveFoundDiv = true
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if haveFoundDiv {
				break
			}
			f(c)
		}
	}
	f(doc)
	return urls, nil
}
