package parser

import (
	"MyProject/house/engine"
	"MyProject/house/model"
	"fmt"
	"regexp"
	"strconv"
)

var (
	data = make([]string, 10)

	houseRe = regexp.MustCompile(`<div class="tt">(.*)</div>`)

	houseMsgRe = regexp.MustCompile(`<span class="rcont">(.*)</span>`)
)

func ParseHouse(bytes []byte, url string) engine.ParseResult {
	fmt.Printf("Get house from %s\n", url)
	house := model.HouseMsg{}
	resultHouse := houseRe.FindAllSubmatch(bytes, -1)
	resultHouseMsg := houseMsgRe.FindAllSubmatch(bytes, -1)
	for _, item := range resultHouse {
		data = append(data, string(item[2]))
	}
	for _, item := range resultHouseMsg {
		data = append(data, string(item[2]))
	}
	fmt.Println(len(data))
	if len(data) != 10 {
		fmt.Println("数据匹配出错,请修改正则表达式")
		return engine.ParseResult{}
	}
	house.House_type = data[0]
	house.Area, _ = strconv.ParseFloat(data[1], 64)
	house.Price, _ = strconv.ParseInt(data[2], 10, 64)
	house.Orientation = data[3]
	house.Floor = data[4]
	house.Decor = data[5]
	house.Is_Elevator = data[6]
	house.Property = data[7]
	house.Structure = data[8]
	house.Year, _ = strconv.ParseInt(data[9], 10, 64)
	result := engine.ParseResult{
		Items: []interface{}{house},
	}
	return result
}
