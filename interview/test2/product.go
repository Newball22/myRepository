package main

type Product interface {
	SetName(name string)
	GetName() string
}

type Product1 struct {
	name string
}

func (p1 *Product1) SetName(name string) {
	p1.name = name
}

func (p1 *Product1) GetName() string {
	return "产品1的名字:" + p1.name
}

type Product2 struct {
	name string
}

func (p2 *Product2) SetName(name string) {
	p2.name = name
}

func (p2 *Product2) GetName() string {
	return "产品2的名字:" + p2.name
}

type productType int

const (
	p1 productType = iota
	p2
)

//实现简单的工厂类
type productFactory struct {
}

func (pf *productFactory) Create(pt productType) Product {
	if pt == p1 {
		return &Product1{}
	}

	if pt == p2 {
		return &Product2{}
	}
	return nil
}
