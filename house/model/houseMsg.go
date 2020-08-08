package model

type HouseMsg struct {
	//户型
	House_type string
	//建筑面积
	Area float64
	//每平方单价(元）
	Price int64
	//房子朝向
	Orientation string
	//楼层高低
	Floor string
	//装修
	Decor string
	//是否有电梯
	Is_Elevator string
	//产权性质「商品房或小产权」
	Property string
	//建筑结构
	Structure string
	//建筑年代
	Year int64
}
