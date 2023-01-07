package controller

import "fmt"

func HandleAdd(obj interface{}) {
	fmt.Println("Handle Add func")
}

func HanldeDelete(obj interface{}) {
	fmt.Println("Handle delete func")
}

func HandleUpdate(oldobj interface{}, obj interface{}) {
	fmt.Println("Handle update func")
}
