package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const buttonPin = 17 // ボタンの GPIO 番号

func Mainrecord() {
	// GPIO 初期化
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}
	defer rpio.Close()

	// GPIO ピンを入力モードに設定
	button := rpio.Pin(buttonPin)
	button.Input()
	button.PullUp() // プルアップ抵抗を有効化（ボタンが押されると LOW になる）

	var cmd *exec.Cmd
	recording := false

	for {
		if button.Read() == rpio.Low { // ボタンが押されたら録音開始
			/*if !recording {
				log.Println("Recording started...")
				cmd = exec.Command("arecord", "-D", "plughw:1,0", "-f", "cd", "-t", "wav", "output.wav")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Start()
				if err != nil {
					log.Fatal(err)
				}
				recording = true
			}*/
			fmt.Println("Recording started...")
		} else { // ボタンが離されたら録音停止
			/*if recording {
				log.Println("Recording stopped.")
				err := cmd.Process.Kill() // `arecord` を強制終了
				if err != nil {
					log.Fatal(err)
				}
				recording = false
			}*/
			fmt.Println("Recording stopped...")
		}
		time.Sleep(100 * time.Millisecond) // CPU 使用率を抑えるために待機
	}
}
