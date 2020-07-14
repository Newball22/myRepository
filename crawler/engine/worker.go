package engine

import (
	"crawler/fetcher"
	"log"
)

func worker(request Request, i int) (ParseResult, error) {
	log.Printf("第%dworker工作，Fetching %s\n", i, request.Url)
	content, err := fetcher.Fetch(request.Url)
	if err != nil {
		log.Printf("Fetching error,Url:%s  Error:%v\n", request.Url, err)
		return ParseResult{}, err
	}

	return request.ParseFunc(content), nil
}
