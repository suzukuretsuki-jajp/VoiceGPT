package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"bytes"
	"encoding/json"

	//"fmt"
	//"io/ioutil"
	"net/http"
	"strings"

	speech "cloud.google.com/go/speech/apiv1"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"

	//"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"google.golang.org/api/option"

	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

const (
	apiKeyFile   = "apikey.txt"   // APIキーが保存されたファイル
	questionFile = "question.txt" // 質問を保存するファイル
	answerFile   = "answer.txt"   // 回答を保存するファイル
	geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent"
)

// Gemini APIのレスポンス構造体
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// ファイルからテキストを読み込む関数
func loadTextFromFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// ファイルにテキストを書き込む関数
func saveTextToFile(filePath, text string) error {
	return ioutil.WriteFile(filePath, []byte(text), 0644)
}

// Gemini API に質問を送る関数
func askGemini(apiKey, question string) (string, error) {
	// リクエストボディの作成
	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": question},
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	// HTTPリクエストを作成
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?key=%s", geminiAPIURL, apiKey), bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// HTTPリクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// レスポンスを読み取る
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// JSONを解析
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", err
	}

	// 回答を取得
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "No response from Gemini", nil
}

func Mainfunc() {
	// サービスアカウントキーのパス
	serviceAccountKeyPath := "C:/Users/sakur/tmciteeep-230010-voicegpt-0f968dbeffbc.json"

	s2taudioFilePath := "./whatyourname.wav" // Speech-to-Textで入力するWAVファイルのパス
	s2toutputFilePath := "./question.txt"    // Speech-to-Textから出力するテキストファイルのパス

	t2sInputFilePath := "./answer.txt"   // Text-to-Speechで入力するテキストファイルのパス
	t2sOutputAudioPath := "./output.wav" // Text-to-Speechから出力するWAVファイルのパス

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

	// APIキーをファイルから読み込む
	apiKey, err := loadTextFromFile(apiKeyFile)
	if err != nil {
		fmt.Println("Error: Failed to load API key:", err)
		return
	}

	// 質問をファイルから読み込む
	question, err := loadTextFromFile(questionFile)
	if err != nil {
		fmt.Println("Error: Failed to load question:", err)
		return
	}

	// Gemini API に問い合わせ
	answer, err := askGemini(apiKey, question)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 回答をファイルに保存
	if err := saveTextToFile(answerFile, answer); err != nil {
		fmt.Println("Error: Failed to save answer:", err)
		return
	}

	fmt.Println("質問:", question)
	fmt.Println("回答を", answerFile, "に保存しました。")

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
