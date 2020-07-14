package engine

type Request struct {
	Url       string                   //请求的地址
	ParseFunc func([]byte) ParseResult //解析函数
}

type ParseResult struct {
	Requests []Request //解析出的请求
	Items    []Item    //解析出的内容，是一个空接口类型，意味着具体的数据结构可以由自己定义{第一阶段}
}

type Item struct {
	Url     string      // 个人信息Url地址
	Type    string      // table
	Id      string      // Id
	Payload interface{} // 详细信息
}

/*
//测试用的
func NilParseFun([]byte) ParseResult {
	return ParseResult{}
}
*/
