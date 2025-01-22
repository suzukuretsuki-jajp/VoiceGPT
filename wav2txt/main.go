package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	//"fmt"
	//"io/ioutil"
	//"log"

	//"os"

	speech "cloud.google.com/go/speech/apiv1"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"

	//"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	"google.golang.org/api/option"

	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

// GeminiRequest represents the request payload for the Gemini API.
type GeminiRequest struct {
	Contents []struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

// GeminiResponse represents the response payload from the Gemini API.
type GeminiResponse struct {
	CachedContent string `json:"cachedContent"`
	Contents      []struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

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

	geminiInputFilePath := "./geminiinput.txt"
	geminiOutputFilePath := "./geminioutput.txt"
	geminiAPIBaseURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key="

	// サービスアカウントキーを読み込む
	keyFileData, err := ioutil.ReadFile(serviceAccountKeyPath)
	if err != nil {
		log.Fatalf("Failed to read service account key file: %v", err)
	}

	// JSONデータを解析してAPIキーを抽出する
	var keyData map[string]interface{}
	err = json.Unmarshal(keyFileData, &keyData)
	if err != nil {
		log.Fatalf("Failed to unmarshal key file: %v", err)
	}

	apiKey, ok := keyData["api_key"].(string) // APIキーのフィールド名を正確に設定
	if !ok || apiKey == "" {
		log.Fatalf("API key not found in the key file")
	}

	// 完全なAPIエンドポイントURL
	geminiAPIURL := geminiAPIBaseURL + apiKey

	// Gemini APIの入力テキストを読み込む
	geminiInputText, err := ioutil.ReadFile(geminiInputFilePath)
	if err != nil {
		log.Fatalf("Failed to read Gemini input text file: %v", err)
	}

	// Gemini APIリクエストの作成
	geminiReq := GeminiRequest{
		Contents: []struct {
			Role  string `json:"role"`
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Role: "user",
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: string(geminiInputText)},
				},
			},
		},
	}

	// リクエストをJSONにシリアライズ
	reqBody, err := json.Marshal(geminiReq)
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	// HTTPリクエストの作成と送信
	req, err := http.NewRequest("POST", geminiAPIURL, strings.NewReader(string(reqBody)))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスを読み取り
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	// レスポンスの解析
	var geminiResp GeminiResponse
	err = json.Unmarshal(respBody, &geminiResp)
	if err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	// 処理されたテキストを新しいテキストファイルに保存
	processedText := geminiResp.Contents[0].Parts[0].Text
	err = ioutil.WriteFile(geminiOutputFilePath, []byte(processedText), 0644)
	if err != nil {
		log.Fatalf("Failed to write Gemini output text file: %v", err)
	}

	fmt.Printf("Processed text has been saved to %s\n", geminiOutputFilePath)

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
