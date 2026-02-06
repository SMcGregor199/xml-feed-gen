package main

import "fmt"

type ghostType struct{
	id int
	name string
	form string
	types []string
	category string
	dexEntry string
}

var gengar = ghostType{
	id: 94,
	name: "Gengar",
	form: "Normal",
	types: []string{"Ghost", "Poison"},
	category: "Shadow",
	dexEntry: "The shadow of a wicked creature, Gengar is said to steal souls and trap them in dolls.",
}


func main() {
	
	fmt.Println(blogDataEndPoint)
}