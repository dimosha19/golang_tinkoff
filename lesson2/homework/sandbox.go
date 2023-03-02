package main

import "fmt"

func main() {
	//p := tagcloud.New()
	//p.AddTag("1")
	//p.AddTag("1")
	//p.AddTag("1")
	//fmt.Println(p.dict["1"])

	p := map[string]int{
		"1": 11,
		"2": 22,
	}
	for i, j := range p {
		fmt.Println(i, j)
	}
}
