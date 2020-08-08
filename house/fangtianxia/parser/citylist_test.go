package parser

import (
	"MyProject/house/fetcher"
	"fmt"
	"testing"
)

//const cityListRe = `<a class="red" href="(//[a-z].esf.fang.com)"\w+>([^<]+)</a>`

func TestParseCityList(t *testing.T) {
	url := "https://tj.esf.fang.com/newsecond/esfcities.aspx"
	data, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println("Fetch failed err:", err)
	}
	result := ParseCityList(data, url)
	fmt.Println(string(data))
	fmt.Println(result)
}
