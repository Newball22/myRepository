package parser

import (
	"MyProject/house/engine"
	"fmt"
	"regexp"
)

const cityListRe = `<a class="red" href="(//[A-Za-z]+.esf.fang.com)".*>([^<]+)</a>`

func ParseCityList(byte []byte, url string) engine.ParseResult {
	fmt.Printf("Get data from %s\n", url)
	re := regexp.MustCompile(cityListRe)
	//subMatch是[][][]byte类型数据
	//第一个[]表示匹配到多少条数据，第二个表示匹配的数据中要提取的内容
	subMatch := re.FindAllSubmatch(byte, -1)
	result := engine.ParseResult{}
	for _, item := range subMatch {
		result.Items = append(result.Items, "City:"+string(item[2]))
		result.Request = append(result.Request, engine.Request{
			Url:       "http:" + string(item[1]),
			ParseFunc: ParseCity,
		})
	}
	return result
}
