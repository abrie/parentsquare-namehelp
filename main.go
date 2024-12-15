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

func getSessionData() (string, map[string]string, error) {
	resp, err := http.Get("https://www.parentsquare.com/sessions")
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	re := regexp.MustCompile(`<form[^>]*action="/sessions"[^>]*>.*?<input[^>]*name="authenticity_token"[^>]*value="([^"]+)"`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", nil, fmt.Errorf("authenticity token not found")
	}
	authenticityToken := matches[1]

	cookie := resp.Header.Get("Set-Cookie")
	if cookie == "" {
		return "", nil, fmt.Errorf("cookie not found")
	}

	psCookies := extractPsCookies(resp.Cookies())

	return authenticityToken, psCookies, nil
}

func extractPsCookies(cookies []*http.Cookie) map[string]string {
	psCookies := make(map[string]string)
	for _, cookie := range cookies {
		if strings.HasPrefix(cookie.Name, "ps_") {
			psCookies[cookie.Name] = cookie.Value
		}
	}
	return psCookies
}

func login(authenticityToken, username, password string, cookies map[string]string) (map[string]string, error) {
	data := url.Values{}
	data.Set("utf8", "✓")
	data.Set("authenticity_token", authenticityToken)
	data.Set("session[email]", username)
	data.Set("session[password]", password)
	data.Set("commit", "Sign In")

	req, err := http.NewRequest("POST", "https://www.parentsquare.com/sessions", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for name, value := range cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Log the redirect
			fmt.Printf("Redirected to: %s\n", req.URL)
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound {
		fmt.Println("Login successful with redirect")
		psCookies := extractPsCookies(resp.Cookies())
		return psCookies, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status code: %d", resp.StatusCode)
	}

	psCookies := extractPsCookies(resp.Cookies())
	return psCookies, nil
}

func queryAutocompleteService(schoolID, limit, chat, query string, cookies map[string]string) (string, error) {
	url := fmt.Sprintf("https://www.parentsquare.com/schools/%s/users/autocomplete?limit=%s&chat=%s&query=%s", schoolID, limit, chat, query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("accept-language", "en")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://www.parentsquare.com/schools/732/users/24399867/chats/new?private=true")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	for name, value := range cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("query failed with status code: %d", resp.StatusCode)
	}

	return string(body), nil
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

	authenticityToken, cookies, err := getSessionData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Authenticity Token:", authenticityToken)
	fmt.Println("Cookies:", cookies)

	psCookies, err := login(authenticityToken, username, password, cookies)
	if err != nil {
		fmt.Println("Login Error:", err)
		return
	}
	fmt.Println("PS Cookies:", psCookies)

	autocompleteResults, err := queryAutocompleteService("732", "25", "1", "cha", psCookies)
	if err != nil {
		fmt.Println("Autocomplete Query Error:", err)
		return
	}
	fmt.Println("Autocomplete Results:", autocompleteResults)
}
