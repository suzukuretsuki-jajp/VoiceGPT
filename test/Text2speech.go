package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func Text2speech() {
	// 認証情報ファイルのパス
	credentialsFile := "/home/ubuntu/voicegpt/VoiceGPT/test/tmciteeep-230010-voicegpt-0f968dbeffbc.json"

	// 入出力ファイルのパス（Text-to-Speech用）
	ttsInputTextFile := "answer.txt"   // 入力テキストファイル
	ttsOutputAudioFile := "output.wav" // 出力音声ファイル

	// Google Cloud Text-to-Speechクライアントを作成
	ctx := context.Background()
	client, err := texttospeech.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// テキストファイルを読み込む
	textData, err := ioutil.ReadFile(ttsInputTextFile)
	if err != nil {
		log.Fatalf("Failed to read text file: %v", err)
	}

	// Text-to-Speechのリクエストを作成
	req := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: string(textData)},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "ja-JP",           // 日本語
			Name:         "ja-JP-Wavenet-A", // ボイスの選択（適宜変更してください）
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16, // 出力音声のエンコーディング
		},
	}

	// 音声合成を実行
	resp, err := client.SynthesizeSpeech(ctx, req)
	if err != nil {
		log.Fatalf("Failed to synthesize speech: %v", err)
	}

	// 結果をWAVファイルに書き込む
	err = ioutil.WriteFile(ttsOutputAudioFile, resp.AudioContent, 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Printf("Audio content written to file: %s\n", ttsOutputAudioFile)
}
