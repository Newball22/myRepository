package engine

//func Run(seeds ...Request) {
//	//建立任务队列
//	var requests []Request
//	for _, r := range seeds {
//		requests = append(requests, r)
//	}
//	for len(requests) > 0 {
//		request := requests[0]
//		requests = requests[1:]
//		url := request.Url
//		fmt.Printf("Fetching %s\n", url)
//		content, err := fetcher.Fetch(url)
//		if err != nil {
//			fmt.Printf("Fetch failed, Url:%s err:%v\n", url, err)
//			continue
//		}
//		parseResult := request.ParseFunc(content, url)
//		requests = append(requests, parseResult.Request...)
//	}
//}
