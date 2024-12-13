package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func getAuthenticityToken() (string, error) {
	resp, err := http.Get("https://www.parentsquare.com/signin")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`action="/sessions".*?name="authenticity_token" value="([^"]+)"`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", fmt.Errorf("authenticity token not found")
	}

	return matches[1], nil
}

func main() {
	token, err := getAuthenticityToken()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Authenticity Token:", token)
}
