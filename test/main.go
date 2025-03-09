package main

import (
	"fmt"
)

func main() {
	fmt.Println("welcome to gemini world")
	for {
		fmt.Println("push to start rec")
		Mainrecord()
		fmt.Println("run Speech-to-Text...")
		Speech2text()
		fmt.Println("finished Speech-to-Text")
		fmt.Println("wakeup gemini API...")
		Maingemini()
		fmt.Println("run Text-to-Speech...")
		Text2speech()
		Mainspeak()
	}
}
