package engine

type ParseResult struct {
	Request []Request
	Items   []interface{}
}

type Request struct {
	Url       string
	ParseFunc func([]byte, string) ParseResult
}

//测试用
func NilParseFun([]byte) ParseResult {
	return ParseResult{}
}
