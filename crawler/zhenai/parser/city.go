package parser

import (
	"crawler/engine"
	"regexp"
)

//(每个用户的url)+(用户名)
var cityRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)

// 用户性别正则，因为在用户详情页没有性别信息，所以性别直接在用户列表页面获取
var sexRe = regexp.MustCompile(`<td width="180"><span class="grayL">性别：</span>([^<]+)</td>`)

//下一页
var nextPageUrlRe = regexp.MustCompile(
	`<li class="paging-item"><a href="(http://www.zhenai.com/zhenghun/[^"]+)">下一页</a>`)

//每个城市页面用户解析器
func ParseCity(bytes []byte) engine.ParseResult {
	subMatch := cityRe.FindAllSubmatch(bytes, -1)
	genderMatch := sexRe.FindAllSubmatch(bytes, -1)

	result := engine.ParseResult{} //定义一个空的ParseResult，方便接下来存储数据使用
	for k, item := range subMatch {
		url := string(item[1])
		name := string(item[2])
		gender := string(genderMatch[k][1])

		//result.Items = append(result.Items, "User:"+name)
		result.Requests = append(result.Requests, engine.Request{
			Url: url,
			ParseFunc: func(bytes []byte) engine.ParseResult {
				return ParseProfile(bytes, name, gender, url)

			},
		})

	}
	//查找下一页
	findSubMatch := nextPageUrlRe.FindAllSubmatch(bytes, -1)
	for _, m := range findSubMatch {
		result.Requests = append(result.Requests, engine.Request{
			Url:       string(m[1]),
			ParseFunc: ParseCity,
		})
	}
	return result
}
