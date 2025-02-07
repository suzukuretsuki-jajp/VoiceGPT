package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	audioFile = "/home/pi/audio/output.wav" // 🔹 再生する音声ファイルのパス
)

func Mainspeak() {
	// Raspberry Pi のオーディオ出力を 3.5mm ジャックに設定
	fmt.Println("Setting audio output to 3.5mm jack...")
	err := exec.Command("amixer", "cset", "numid=3", "1").Run()
	if err != nil {
		fmt.Println("Failed to set audio output:", err)
		return
	}
	fmt.Println("Audio output set to 3.5mm jack.")

	// 音声ファイルの存在を確認
	if _, err := os.Stat(audioFile); os.IsNotExist(err) {
		fmt.Println("Error: Audio file not found:", audioFile)
		return
	}

	// SoX の `play` コマンドで音声を再生
	fmt.Println("Playing audio:", audioFile)
	cmd := exec.Command("play", audioFile)
	err = cmd.Start()
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

	// 再生後に短い待機を入れる（音の途切れ対策）
	time.Sleep(1 * time.Second)
}
