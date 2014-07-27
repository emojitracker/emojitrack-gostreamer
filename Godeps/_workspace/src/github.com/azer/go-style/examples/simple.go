package main

import (
	. "github.com/azer/go-style"
	"fmt"
)

func main () {
	fmt.Println(Style("bold red", "\n Bold red "))
	fmt.Println(Style("yellow greenBg", " yellow greenBg "))
	fmt.Println(Style(".blueBg .white .bold", " blueBg white bold \n"))
}
