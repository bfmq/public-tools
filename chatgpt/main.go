package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type chatRequest struct {
	Model     string `json:"model"`
	MaxTokens int    `json:"max_tokens"`
	Prompt    string `json:"prompt"`
}

type chatResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

var (
	// API Point
	endPoint = "https://api.openai.com/v1/"
	apiKey   = ""

	msgFlag string
)

func init() {
	flag.StringVar(&msgFlag, "msg", "", "文本消息")
}

func main() {
	flag.Parse()

	if msgFlag == "" {
		fmt.Println("请输入文本信息")
		fmt.Println("chat_gpt_demo -msg 文本信息")
		return
	}

	cq := &chatRequest{
		Model:     "gpt-3.5-turbo",
		MaxTokens: 500,
		Prompt:    msgFlag,
	}

	cr, err := request(cq)
	if err != nil {
		fmt.Println("错误信息：", err)
	}

	if cr != nil && len(cr.Choices) > 0 {
		fmt.Printf("ChatGPT：")
		fmt.Printf(cr.Choices[0].Text)
		fmt.Println("")
	}

}

func request(cq *chatRequest) (cr *chatResponse, err error) {

	if cq == nil && cq.Prompt == "" {
		err = errors.New("需要文本信息")
		return
	}

	b, err := json.Marshal(&cq)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", endPoint, bytes.NewReader(b))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("User-Agent", "chatGPT/1 CFNetwork/1402.0.8 Darwin/22.2.0")
	req.Header.Add("Accept-encoding", "gzip, deflate, br")
	req.Header.Add("Accept-language", "zh-CN,zh-Hans;q=0.9")
	client := &http.Client{
		Timeout:   300 * time.Second,
		Transport: http.DefaultTransport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	cr = &chatResponse{}

	err = json.Unmarshal(b, &cr)
	if err != nil {
		return
	}
	return
}
