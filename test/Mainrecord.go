package main

import (
	"fmt"
	"log"
	"os"
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
	nowrecord := false
	recording := false

	for {
		if button.Read() == rpio.Low { // ボタンが押されたら録音開始
			if !recording {
				log.Println("Recording started...")
				cmd = exec.Command("arecord", "-D", "plughw:2,0", "-f", "S16_LE", "-r", "48000", "-c", "1", "-t", "wav", "input.wav")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Start()
				if err != nil {
					log.Fatal(err)
				}
				recording = true
				nowrecord = true

			}
			fmt.Println("Recording started...")
		} else { // ボタンが離されたら録音停止
			if recording {
				log.Println("not Recording now")
				err := cmd.Process.Kill() // `arecord` を強制終了
				if err != nil {
					log.Fatal(err)
				}
				recording = false
			}
			if nowrecord {
				fmt.Println("Recording stopped")
				break
			}

		}
		time.Sleep(100 * time.Millisecond) // CPU 使用率を抑えるために待機
	}
}
