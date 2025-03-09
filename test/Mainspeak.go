package main

import (
	"fmt"
	//"os"
	"os/exec"
	"time"
)

const (
	audioFile = "/home/ubuntu/voicegpt/VoiceGPT/test/answer.wav" // 🔹 再生する音声ファイルのパス
)

func main() {

	// SoX の `play` コマンドで音声を再生
	fmt.Println("Playing audio:", audioFile)
	cmd := exec.Command("play", audioFile)
	err = cmd.Start()
	if err != nil {
		fmt.Println("Failed to start playback:", err)
		return
	}

	// 再生完了を待機
	err := cmd.Wait()
	if err != nil {
		fmt.Println("Error during playback:", err)
	} else {
		fmt.Println("Playback finished.")
	}

	// 再生後に短い待機を入れる（音の途切れ対策）
	time.Sleep(1 * time.Second)
}
