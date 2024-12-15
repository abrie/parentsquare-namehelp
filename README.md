# parentsquare-namehelp

## Login Method

This repository includes a method to log in to ParentSquare using an authenticity token, username, password, and cookie.

### Usage

1. Retrieve the authenticity token and cookie from the sessions page.
2. Call the `login` method with the authenticity token, username, password, and cookie.
3. Parse the username and password from a JSON file specified by a command line argument.

### Example JSON File

```json
{
  "login": {
    "username": "username_here",
    "password": "password_here"
  }
}
```

## Autocomplete Query Method

This repository now includes a method to query the autocomplete service using the provided cURL command as a model.

### Usage

1. Call the `queryAutocomplete` method with the cookie, CSRF token, and query string.
2. The method constructs the request using the provided cURL command as a model.
3. The method sends the request and returns the response.

### Example Query

```go
query := "cha" // Example query
response, err := queryAutocomplete(cookie, csrfToken, query)
if err != nil {
    fmt.Println("Autocomplete Query Error:", err)
    return
}
fmt.Println("Autocomplete Response:", response)
```
