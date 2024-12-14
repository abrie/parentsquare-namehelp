package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Credentials struct {
	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"login"`
}

func parseCredentials(filename string) (string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	var creds Credentials
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&creds)
	if err != nil {
		return "", "", err
	}

	return creds.Login.Username, creds.Login.Password, nil
}

func getSessionData() (string, string, error) {
	resp, err := http.Get("https://www.parentsquare.com/sessions")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	re := regexp.MustCompile(`<form[^>]*action="/sessions"[^>]*>.*?<input[^>]*name="authenticity_token"[^>]*value="([^"]+)"`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", "", fmt.Errorf("authenticity token not found")
	}
	authenticityToken := matches[1]

	cookie := resp.Header.Get("Set-Cookie")
	if cookie == "" {
		return "", "", fmt.Errorf("cookie not found")
	}

	return authenticityToken, cookie, nil
}

func login(authenticityToken, username, password, cookie string) error {
	data := url.Values{}
	data.Set("utf8", "✓")
	data.Set("authenticity_token", authenticityToken)
	data.Set("session[email]", username)
	data.Set("session[password]", password)
	data.Set("commit", "Sign In")

	req, err := http.NewRequest("POST", "https://www.parentsquare.com/sessions", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Origin", "https://www.parentsquare.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://www.parentsquare.com/signin")
	req.Header.Set("Sec-CH-UA", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Log the response code
	fmt.Printf("Response Code: %d\n", resp.StatusCode)

	// Show the contents of the response's 'Location' header, if present
	location := resp.Header.Get("Location")
	if location != "" {
		fmt.Printf("Location Header: %s\n", location)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_json_file>")
		return
	}

	jsonFile := os.Args[1]
	username, password, err := parseCredentials(jsonFile)
	if err != nil {
		fmt.Println("Error parsing credentials:", err)
		return
	}

	authenticityToken, cookie, err := getSessionData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Authenticity Token:", authenticityToken)
	fmt.Println("Cookie:", cookie)

	err = login(authenticityToken, username, password, cookie)
	if err != nil {
		fmt.Println("Login Error:", err)
		return
	}
	fmt.Println("Login successful")
}
