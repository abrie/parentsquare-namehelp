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

This repository now includes a method to query the ParentSquare autocomplete service.

### Usage

1. Call the `queryAutocompleteService` method with the school ID, limit, chat, and query parameters.
2. The method returns the response from the autocomplete service.

### Example

```go
autocompleteResults, err := queryAutocompleteService("732", "25", "1", "cha", psCookies)
if err != nil {
    fmt.Println("Autocomplete Query Error:", err)
    return
}
fmt.Println("Autocomplete Results:", autocompleteResults)
```

## HTTP Server

This repository now includes an HTTP server that provides a route to the autocomplete method.

### Running the Server

1. Ensure you have a JSON file with your ParentSquare credentials (username and password).
2. Run the server using the following command:

```sh
go run main.go <path_to_json_file>
```

### Example Usage of the `/autocomplete` Route

1. Start the server as described above.
2. Make a GET request to the `/autocomplete` route with the required query parameters (`school_id`, `limit`, `chat`, `query`).

Example:

```sh
curl "http://localhost:8080/autocomplete?school_id=732&limit=25&chat=1&query=cha"
```
