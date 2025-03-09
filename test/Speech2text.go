package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	"google.golang.org/api/option"

	//"google.golang.org/genproto/googleapis/cloud/speech/v1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func Speech2text() {
	// 認証情報ファイルのパス
	credentialsFile := "/home/ubuntu/voicegpt/VoiceGPT/test/tmciteeep-230010-voicegpt-0f968dbeffbc.json"

	// 入出力ファイルのパス（Speech-to-Text用）
	sttInputAudioFile := "input.wav"    // 入力音声ファイル
	sttOutputTextFile := "question.txt" // 出力テキストファイル

	// Google Cloud Speech-to-Textクライアントを作成
	ctx := context.Background()
	client, err := speech.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// 音声ファイルを読み込む
	audioData, err := ioutil.ReadFile(sttInputAudioFile)
	if err != nil {
		log.Fatalf("Failed to read audio file: %v", err)
	}

	// 音声認識のリクエストを作成
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 48000,   // サンプルレート（適宜変更してください）
			LanguageCode:    "ja-JP", // 日本語
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: audioData},
		},
	}

	// 音声認識を実行
	resp, err := client.Recognize(ctx, req)
	if err != nil {
		log.Fatalf("Failed to recognize: %v", err)
	}

	// 結果をTXTファイルに書き込む
	var transcription string
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			transcription += alt.Transcript + "\n"
		}
	}

	err = ioutil.WriteFile(sttOutputTextFile, []byte(transcription), 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Printf("Transcription saved to %s\n", sttOutputTextFile)
}
