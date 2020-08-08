package parser

import (
	"MyProject/house/engine"
	"fmt"
	"regexp"
)

const cityRe = `<a.*href="(\/chushou\/\d_\d+.htm).*target.*>`

func ParseCity(bytes []byte, url string) engine.ParseResult {
	fmt.Printf("Get city from %s\n", url)
	re := regexp.MustCompile(cityRe)
	subMatch := re.FindAllSubmatch(bytes, -1)

	result := engine.ParseResult{}
	for _, item := range subMatch {
		result.Request = append(result.Request, engine.Request{
			Url:       url + string(item[1]),
			ParseFunc: ParseHouse,
		})
	}
	return result
}
