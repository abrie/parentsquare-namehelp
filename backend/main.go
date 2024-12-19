package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

type Credentials struct {
	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"login"`
}

type Config struct {
	Autocomplete struct {
		SchoolID string `json:"school_id"`
		Limit    string `json:"limit"`
		Chat     string `json:"chat"`
	} `json:"autocomplete"`
}

type Server struct {
	psCookies map[string]string
	config    Config
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

func parseConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
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

func (s *Server) queryAutocompleteService(schoolID, limit, chat, query string) (string, error) {
	fmt.Printf("Starting autocomplete query for: %s\n", query) // Pff3e
	url := fmt.Sprintf("https://www.parentsquare.com/schools/%s/users/autocomplete?limit=%s&chat=%s&query=%s", schoolID, limit, chat, query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("accept", "application/json")

	for name, value := range s.psCookies {
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

	fmt.Printf("Completed autocomplete query for: %s\n", query) // Pff3e
	return string(body), nil
}

func (s *Server) autocompleteHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received autocomplete query: %s\n", query) // P2687
	results, err := s.queryAutocompleteService(s.config.Autocomplete.SchoolID, s.config.Autocomplete.Limit, s.config.Autocomplete.Chat, query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Autocomplete Query Error: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Autocomplete results for query '%s': %s\n", query, results) // P2687
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(results))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <path_to_credentials_json_file> <path_to_config_json_file>")
		return
	}

	credentialsFile := os.Args[1]
	configFile := os.Args[2]

	username, password, err := parseCredentials(credentialsFile)
	if err != nil {
		fmt.Println("Error parsing credentials:", err)
		return
	}

	config, err := parseConfig(configFile)
	if err != nil {
		fmt.Println("Error parsing config:", err)
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

	server := &Server{psCookies: psCookies, config: config}

	http.HandleFunc("/api/autocomplete", server.autocompleteHandler)

	srv := &http.Server{
		Addr: ":8080",
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe(): %v\n", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown Failed:%+v", err)
	}
	fmt.Println("Server Exited Properly")
}
