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
	"net/http/httputil"
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

func logCurlEquivalent(url, method, data, cookie string, headers map[string]string) {
	var curlCmd strings.Builder
	curlCmd.WriteString("curl -X ")
	curlCmd.WriteString(method)
	curlCmd.WriteString(" '")
	curlCmd.WriteString(url)
	curlCmd.WriteString("'")

	for key, value := range headers {
		curlCmd.WriteString(" -H '")
		curlCmd.WriteString(fmt.Sprintf("%s: %s", key, value))
		curlCmd.WriteString("'")
	}

	if data != "" {
		curlCmd.WriteString(" --data '")
		curlCmd.WriteString(data)
		curlCmd.WriteString("'")
	}

	if cookie != "" {
		curlCmd.WriteString(" --cookie '")
		curlCmd.WriteString(cookie)
		curlCmd.WriteString("'")
	}

	fmt.Println(curlCmd.String())
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
	req.Header.Set("Cookie", cookie)

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Cookie":       cookie,
	}

	// Log the 'curl' CLI equivalent of the request
	logCurlEquivalent("https://www.parentsquare.com/sessions", "POST", data.Encode(), cookie, headers)

	// Dump the request
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}
	fmt.Printf("Request:\n%s\n", string(requestDump))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Dump the response
	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return err
	}
	fmt.Printf("Response:\n%s\n", string(responseDump))

	// Log the response code
	fmt.Printf("Response Code: %d\n", resp.StatusCode)

	// Show the contents of the response's 'Location' header, if present
	location := resp.Header.Get("Location")
	if location != "" {
		fmt.Printf("Location Header: %s\n", location)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Print the length of response text
	fmt.Printf("Response Text: %d\n", len(string(body)))

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
}
