package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name string
	Age  int
}

func main() {
	typ := reflect.TypeOf(User{})
	fmt.Println(typ.Name(), typ.Kind())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fmt.Println(field.Name)
	}

	// 使用指针的情况
	typ = reflect.Indirect(reflect.ValueOf(&User{})).Type()
	fmt.Println(typ.Name())

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fmt.Println(field.Name)
	}
}
