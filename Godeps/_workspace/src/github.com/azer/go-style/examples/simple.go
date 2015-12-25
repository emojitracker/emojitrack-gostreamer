package main

import (
	"fmt"
	. "github.com/mroth/emojitrack-gostreamer/Godeps/_workspace/src/github.com/azer/go-style"
)

func main() {
	fmt.Println(Style("bold red", "\n Bold red "))
	fmt.Println(Style("yellow greenBg", " yellow greenBg "))
	fmt.Println(Style(".blueBg .white .bold", " blueBg white bold \n"))
}
