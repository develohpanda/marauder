package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	rp := validateArguments()
	readFirstPage(rp)
}

func readFirstPage(rp runParams) {
	baseURL := fmt.Sprintf("https://graph.facebook.com/v2.3/%s/messages", rp.threadID)
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		log.Fatal("Could not create request")
		return
	}

	q := req.URL.Query()
	q.Add("access_token", rp.accessToken)
	q.Add("limit", "2")

	req.URL.RawQuery = q.Encode()
	client := http.Client{}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Fatal("Could not load data from graph api")
		return
	}

	fmt.Println()
	fmt.Println(processResponse(resp))
}

func processResponse(resp *http.Response) string {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
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
