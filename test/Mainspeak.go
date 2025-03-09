package main

import (
	"fmt"
	"os/exec"
)

func Mainspeak() {
	// SoX の `play` コマンドで音声を再生
	fmt.Println("Playing audio:", "answer.wav")
	cmd := exec.Command("play", "answer.wav")

	// 再生開始
	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to start playback:", err)
		return
	}

	// 再生完了を待機
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error during playback:", err)
	} else {
		fmt.Println("Playback finished.")
	}

}
