package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ファイルのパス
const (
	apiKeyFile     = "apikey.txt"     // APIキーが保存されたファイル
	questionFile   = "question.txt"   // 質問を保存するファイル
	answerFile     = "answer.txt"     // 回答を保存するファイル
	geminiAPIURL   = "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent"
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

func main() {
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
}
