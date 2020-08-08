package engine

import (
	"MyProject/house/fetcher"
	"log"
)

func worker(request Request) (ParseResult, error) {
	url := request.Url
	log.Printf("Fetch %s\n", url)
	content, err := fetcher.Fetch(url)
	if err != nil {
		log.Printf("Fetch error, Url:%s %v\n", url, err)
		return ParseResult{}, err
	}
	return request.ParseFunc(content, url), nil
}
