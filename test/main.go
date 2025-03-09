package main

import (
	"fmt"
)

func main() {
	fmt.Println("welcome to gemini world")
	for {
		fmt.Println("push to start rec")
		Mainrecord()
		Speech2text()
		Maingemini()
		Text2speech()
		Mainspeak()
	}
}
