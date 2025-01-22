package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	speech "cloud.google.com/go/speech/apiv1"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"

	//"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"google.golang.org/api/option"

	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func main() {
	// サービスアカウントキーのパス
	serviceAccountKeyPath := "C:/Users/sakur/VoiceGPT/Googlecloudkey/tmciteeep-230010-voicegpt-21fb464420b2.json"

	s2taudioFilePath := "./testaudio3a.wav"  // Speech-to-Textで入力するWAVファイルのパス
	s2toutputFilePath := "./s2ttesttext.txt" // Speech-to-Textから出力するテキストファイルのパス

	t2sInputFilePath := "./s2ttesttext.txt"    // Text-to-Speechで入力するテキストファイルのパス
	t2sOutputAudioPath := "./t2stestaudio.wav" // Text-to-Speechから出力するWAVファイルのパス

	// 音声ファイルを読み込む
	audioData, err := ioutil.ReadFile(s2taudioFilePath)
	if err != nil {
		log.Fatalf("Failed to read audio file: %v", err)
	}

	// Google Speech-to-Text クライアントを作成
	ctx := context.Background()
	client, err := speech.NewClient(ctx, option.WithCredentialsFile(serviceAccountKeyPath))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// APIに送信するリクエストを構築
	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 48000, // WAVファイルのサンプリングレート
			LanguageCode:    "ja-JP",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: audioData},
		},
	}

	// Speech-to-Text APIにリクエストを送信
	resp, err := client.Recognize(ctx, req)
	if err != nil {
		log.Fatalf("Failed to recognize speech: %v", err)
	}

	// テキストファイルに結果を書き込む
	file, err := os.Create(s2toutputFilePath)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	// UTF-8形式で結果を保存
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			_, err := file.WriteString(fmt.Sprintf("%v\n", alt.Transcript))
			if err != nil {
				log.Fatalf("Failed to write to output file: %v", err)
			}
		}
	}

	fmt.Printf("Transcription has been saved to %s\n", s2toutputFilePath)

	// Text-to-Speechの入力テキストを読み込む
	inputText, err := ioutil.ReadFile(t2sInputFilePath)
	if err != nil {
		log.Fatalf("Failed to read input text file: %v", err)
	}

	// Google Text-to-Speech クライアントを作成
	t2sClient, err := texttospeech.NewClient(ctx, option.WithCredentialsFile(serviceAccountKeyPath))
	if err != nil {
		log.Fatalf("Failed to create Text-to-Speech client: %v", err)
	}
	defer t2sClient.Close()

	// APIに送信するリクエストを構築
	t2sReq := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: string(inputText)},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "ja-JP",            // 日本語の音声
			Name:         "ja-JP-Standard-A", // 日本語のスタンダードな音声
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding:   texttospeechpb.AudioEncoding_LINEAR16, // 16bit WAV
			SampleRateHertz: 48000,                                 // 48kHz
		},
	}

	// Text-to-Speech APIにリクエストを送信
	t2sResp, err := t2sClient.SynthesizeSpeech(ctx, t2sReq)
	if err != nil {
		log.Fatalf("Failed to synthesize speech: %v", err)
	}

	// 音声ファイルに結果を書き込む
	err = ioutil.WriteFile(t2sOutputAudioPath, t2sResp.AudioContent, 0644)
	if err != nil {
		log.Fatalf("Failed to write audio file: %v", err)
	}

	fmt.Printf("Audio content has been saved to %s\n", t2sOutputAudioPath)

}
