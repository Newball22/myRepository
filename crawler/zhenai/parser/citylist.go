package parser

import (
	"crawler/engine"
	"regexp"
)

//(http://www.zhenai.com/zhenghun/guangzhou)(广州)
const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

//解析城市列表的数据
func ParseCityList(bytes []byte) engine.ParseResult {
	re := regexp.MustCompile(cityListRe)
	subMatch := re.FindAllSubmatch(bytes, -1) //这个'-1'是指匹配到末尾
	result := engine.ParseResult{}            //定义一个空的ParseResult
	limit := 2
	for _, item := range subMatch {
		//result.Items = append(result.Items, "爬取到的City:"+string(item[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url:       string(item[1]), //每个城市对应的URL
			ParseFunc: ParseCity,       //使用城市解析器
		})
		limit--
		if limit == 0 {
			break
		}
	}
	return result
}
