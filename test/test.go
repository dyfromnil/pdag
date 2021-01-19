package main

import (
	"fmt"
	"github.com/deckarep/golang-set"
)

func main() {

	var aka mapset.Set
	aka = mapset.NewSet()
	aka.Add(34)
	aka.Add(34)
	aka.Add(34)
	aka.Add(64)
	it := aka.Iterator()
	var intlist []int
	for elem := range it.C {
		intlist = append(intlist, elem.(int))
	}
	it.Stop()
	for _, e := range intlist {
		fmt.Println(e)
	}

	requiredClasses := mapset.NewSet()
	requiredClasses.Add("Cooking")
	requiredClasses.Add("English")
	requiredClasses.Add("Math")
	requiredClasses.Add("Biology")
	requiredClasses.Add(574)
	requiredClasses.Remove("Cooking")
	requiredClasses.Remove("aaaaa")

	it = requiredClasses.Iterator()
	for elem := range it.C {
		fmt.Println(elem)
	}
	it.Stop()

}
