package main

import (
	"fmt"
	"os/exec"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

const (
	gpioPin  = "P1_11"               // GPIO 17 (物理ピン11)
	filePath = "/home/pi/record.wav" // 保存するファイルパス
)

func main() {
	// `periph` を初期化
	if _, err := host.Init(); err != nil {
		fmt.Println("Failed to initialize periph:", err)
		return
	}

	// GPIO 17 を入力モードに設定（プルダウン）
	pin := rpi.P1_11
	pin.In(gpio.PullDown, gpio.FallingEdge)

	fmt.Println("Press and hold the button to record audio...")

	recording := false
	var cmd *exec.Cmd

	for {
		//録音部動作開始
		for {
			if pin.Read() == gpio.High { // ボタンが押されたら録音開始
				if !recording {
					fmt.Println("Recording started...")
					cmd = exec.Command("rec", filePath, "rate", "32k") // SoX で録音開始
					err := cmd.Start()
					if err != nil {
						fmt.Println("Failed to start recording:", err)
						continue
					}
					recording = true
				}
			} else { // ボタンを離したら録音停止
				if recording {
					fmt.Println("Recording stopped.")
					err := cmd.Process.Kill() // `rec` を停止
					if err != nil {
						fmt.Println("Failed to stop recording:", err)
					}
					recording = false
				}
			}
			if pin.Read() != gpio.High {
				break //録音部分のループ終了
			}
			time.Sleep(100 * time.Millisecond) // CPU負荷を減らすためスリープ
		}
		//録音部動作終了

		//gemini動作開始
		Mainfunc()
		//Gemini動作終了

		//生成した音声を再生

	}
}
