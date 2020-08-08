package main

import (
	"fmt"
	"testing"
)

func TestProduct_Create(t *testing.T) {
	p1 := &Product1{}
	p1.SetName("p1")
	fmt.Println(p1.GetName())

	p2 := &Product2{}
	p2.SetName("p2")
	fmt.Println(p2.GetName())
}

func TestProductFactory_Create(t *testing.T) {
	proFactory := productFactory{}

	product1 := proFactory.Create(p1)
	product1.SetName("p1")
	fmt.Println(product1.GetName())

	product2 := proFactory.Create(p2)
	product2.SetName("p2")
	fmt.Println(product2.GetName())
}
