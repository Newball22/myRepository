package parser

import (
	"MyProject/house/fetcher"
	"fmt"
	"testing"
)

func TestParseCity(t *testing.T) {
	url := "https://gz.esf.fang.com"
	data, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println("Fetch failed err:", err)
	}
	result := ParseCity(data, url)
	fmt.Println(string(data))
	fmt.Println(result)
}
