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

	req.Header.Set("accept", "application/json")

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

func autocompleteHandler(w http.ResponseWriter, r *http.Request) {
	schoolID := r.URL.Query().Get("school_id")
	limit := r.URL.Query().Get("limit")
	chat := r.URL.Query().Get("chat")
	query := r.URL.Query().Get("query")

	// Assuming psCookies are available globally or through some other means
	autocompleteResults, err := queryAutocompleteService(schoolID, limit, chat, query, psCookies)
	if err != nil {
		http.Error(w, "Autocomplete Query Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(autocompleteResults))
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

	http.HandleFunc("/autocomplete", autocompleteHandler)
	port := ":8080"
	fmt.Println("Server is running on port", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Server Error:", err)
	}
}
