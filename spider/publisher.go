package spider

import "log"

func Publisher(startUrl string) (map[string][]string, error) {
	content, err := Spider(startUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	urls, err := HtmlParseChapter(content)

	return urls, nil
}
