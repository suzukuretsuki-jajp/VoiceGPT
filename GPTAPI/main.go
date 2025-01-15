package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"bytes"
	"encoding/json"
	"net/http"
)

// GPT APIの認証鍵のパス
const apiURL = "https://api.openai.com/v1/chat/completions"

// GPT APIのリクエスト構造体
type GPTRequest struct {
	Model    string       `json:"model"`
	Messages []GPTMessage `json:"messages"`
}

type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GPT APIのレスポンス構造体
type GPTResponse struct {
	Choices []GPTChoice `json:"choices"`
}

type GPTChoice struct {
	Message GPTMessage `json:"message"`
}

func main() {

	// GPT APIのAPIキーをjsonファイルから取得
	apiKey, err := getAPIKey("../wav2txt/gptkey.json")
	if err != nil {
		log.Fatalf("Failed to get API key: %v", err)
	}

	//初期設定(パス等)

	GPTinputFile := "GPTtesttextIN.txt"
	GPToutputFile := "GPTtesttextOUT.txt"

	/* GPT command zone--------------------------------
	---------------------------------------------------
	---------------------------------------------------
	-------------------------------------------------*/

	// GPTInputファイルを読み込み
	GPTinputData, err := ioutil.ReadFile(GPTinputFile)
	if err != nil {
		log.Fatalf("GPTInputファイルの読み込みに失敗しました: %v", err)
	}

	// GPTRequestを構成
	GPTrequestData := &GPTRequest{
		Model: "gpt-4", // 使用するモデルを指定
		Messages: []GPTMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: string(GPTinputData),
			},
		},
	}

	// JSONエンコードしてGPTRequestのリクエストボディを作成
	GPTreqBody, err := json.Marshal(GPTrequestData)
	if err != nil {
		log.Fatalf("GPTRequestのJSONエンコードに失敗しました: %v", err)
	}
	fmt.Println("gptreqbody", GPTreqBody)

	// HTTPリクエストの作成
	GPTreq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(GPTreqBody))
	if err != nil {
		log.Fatalf("GPT APIリクエストの作成に失敗しました: %v", err)
	}

	// APIキーをヘッダに追加
	GPTreq.Header.Set("Content-Type", "application/json")
	GPTreq.Header.Set("Authorization", "Bearer "+apiKey)

	fmt.Println("gptreq", GPTreq)

	// HTTPクライアントを使ってリクエストを送信
	GPTclient := &http.Client{}
	GPTresp, err := GPTclient.Do(GPTreq)
	if err != nil {
		log.Fatalf("GPT APIリクエストの送信に失敗しました: %v", err)
	}
	defer GPTresp.Body.Close()

	// GPTResponseを解析
	var GPTresponse GPTResponse
	if err := json.NewDecoder(GPTresp.Body).Decode(&GPTresponse); err != nil {
		log.Fatalf("GPT APIレスポンスの解析に失敗しました: %v", err)
	}

	fmt.Println("gptresp", GPTresponse)
	// 結果をGPTOutputファイルに書き込む

	if len(GPTresponse.Choices) == 0 {
		log.Fatalf("No choices found in the GPT response")
	}

	GPTresult := GPTresponse.Choices[0].Message.Content
	err = ioutil.WriteFile(GPToutputFile, []byte(GPTresult), 0644)
	if err != nil {
		log.Fatalf("結果のGPTOutputファイル書き込みに失敗しました: %v", err)
	}

	fmt.Println("処理が完了しました。結果は 'output.txt' に保存されました。")

}

func getAPIKey(filename string) (string, error) {
	// ファイルを読み込む
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read API key file: %v", err)
	}

	// JSONをデコードしてapikeyを取得
	var config map[string]string
	if err := json.Unmarshal(data, &config); err != nil {
		return "", fmt.Errorf("failed to parse API key file: %v", err)
	}

	apiKey, exists := config["apikey"]
	if !exists {
		return "", fmt.Errorf("apikey not found in %s", filename)
	}

	return apiKey, nil
}
