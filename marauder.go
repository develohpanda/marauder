package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	rp := validateArguments()
	msgs := getAllMesasges(rp)
	printMetrics(msgs, rp)
}

func printMetrics(msgs []message, rp runParams) {
	var firstPersonKeywordCount int64
	var secondPersonKeywordCount int64
	var totalCharacters int64
	var searchedCharacters int64

	firstPersonKeywordCount = 0
	secondPersonKeywordCount = 0
	totalCharacters = 0
	searchedCharacters = 0

	for i1 := 0; i1 < len(msgs); i1++ {
		msg := msgs[i1]
		totalCharacters += int64(len(msg.Message))
		for i2 := 0; i2 < len(rp.keywords); i2++ {

			kw := rp.keywords[i2]
			count := strings.Count(msg.Message, kw)
			searchedCharacters += int64(count * len(kw))

			if strings.Contains(msg.From.Name, rp.firstPerson) {
				firstPersonKeywordCount += int64(count)
			} else {
				secondPersonKeywordCount += int64(count)
			}
		}
	}

	fmt.Println(fmt.Sprintf("Total messages: %d", len(msgs)))
	fmt.Println()
	fmt.Println("Keyword counts:")
	fmt.Println(fmt.Sprintf("%s: %d", rp.firstPerson, firstPersonKeywordCount))
	fmt.Println(fmt.Sprintf("%s: %d", rp.secondPerson, secondPersonKeywordCount))
	fmt.Println(fmt.Sprintf("Total: %d", firstPersonKeywordCount+secondPersonKeywordCount))
	fmt.Println()
	fmt.Println(fmt.Sprintf("Total characters: %d", totalCharacters))
	fmt.Println(fmt.Sprintf("Characters of joy: %d", searchedCharacters))
	fmt.Println(fmt.Sprintf("Percentage joy: %f", float64(searchedCharacters)/float64(totalCharacters)))
	fmt.Println()
	fmt.Println("Not funny at all... :(")
}

func getAllMesasges(rp runParams) []message {
	fmt.Println()
	fmt.Println("I solemnly swear that I am up to no good...")
	var msgs []message
	decodedResponse := readFirstPage(rp)

	for {
		msgs = append(msgs, decodedResponse.Data...)

		if len(decodedResponse.Paging.Next) <= 0 {
			fmt.Println("...mischief managed.")
			fmt.Println()
			break
		}
		decodedResponse = readRemainingPages(decodedResponse.Paging.Next)
	}
	return msgs
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

	if argumentCount < 6 {
		log.Fatal("Please add access token, thread id, the two people's names, and which words to look for.")
	}

	accessToken := os.Args[1]
	threadID := os.Args[2]
	firstPerson := os.Args[3]
	secondPerson := os.Args[4]
	keywords := os.Args[5:]

	return runParams{accessToken: accessToken, threadID: threadID, keywords: keywords, firstPerson: firstPerson, secondPerson: secondPerson}
}

type runParams struct {
	accessToken  string
	threadID     string
	keywords     []string
	firstPerson  string
	secondPerson string
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
