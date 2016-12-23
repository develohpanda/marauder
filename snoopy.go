package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	rp := validateArguments()

	var msgs []message
	decodedResponse := readFirstPage(rp)

	for {
		msgs = append(msgs, decodedResponse.Data...)

		fmt.Println(len(msgs))
		if len(decodedResponse.Paging.Next) <= 0 {
			fmt.Println("Done")
			break
		}
		decodedResponse = readRemainingPages(decodedResponse.Paging.Next)
	}
}

func readRemainingPages(fullURL string) messagesResponse {
	resp, err := http.Get(fullURL)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(err)
	}
	return decodeJSONResponse(processResponse(resp))
}

func decodeJSONResponse(byt []byte) messagesResponse {
	res := messagesResponse{}
	err := json.Unmarshal(byt, &res)
	if err != nil {
		panic(err)
	}
	return res
}

func readFirstPage(rp runParams) messagesResponse {
	baseURL := fmt.Sprintf("https://graph.facebook.com/v2.3/%s/messages", rp.threadID)
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("access_token", rp.accessToken)
	q.Add("limit", "500")

	req.URL.RawQuery = q.Encode()
	client := http.Client{}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(err)
	}

	return decodeJSONResponse(processResponse(resp))
}

func processResponse(resp *http.Response) []byte {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func validateArguments() runParams {
	argumentCount := len(os.Args)

	if argumentCount < 4 {
		log.Fatal("Please add access token, thread id, and which words to look for.")
	}

	accessToken := os.Args[1]
	threadID := os.Args[2]
	keywords := os.Args[3:]

	return runParams{accessToken: accessToken, threadID: threadID, keywords: keywords}
}

type runParams struct {
	accessToken string
	threadID    string
	keywords    []string
}

type messagesResponse struct {
	Data   []message `json:"data"`
	Paging paging    `json:"paging"`
}

type paging struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

type message struct {
	ID          string `json:"id"`
	CreatedTime string `json:"created_time"`
	Tags        struct {
		Data []struct {
			Name string `json:"name"`
		} `json:"data"`
	} `json:"tags"`
	From struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		ID    string `json:"id"`
	} `json:"from"`
	To struct {
		Data []struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			ID    string `json:"id"`
		} `json:"data"`
	} `json:"to"`
	Message string `json:"message"`
	Shares  struct {
		Data []struct {
			ID          interface{} `json:"id"`
			Link        string      `json:"link"`
			Name        string      `json:"name"`
			Description interface{} `json:"description"`
			Picture     interface{} `json:"picture"`
		} `json:"data"`
	} `json:"shares,omitempty"`
}
