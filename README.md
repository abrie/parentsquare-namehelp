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

### Detailed Logging

The `login` method now includes detailed logging that mirrors Curl's `--trace-ascii` feature. This logging includes the request method, URL, headers, and body, as well as the response status, headers, and body.

### Example Output

```
Request Method: POST
Request URL: https://www.parentsquare.com/sessions
Request Header: Content-Type: application/x-www-form-urlencoded
Request Header: Cookie: <cookie_value>
Request Body: utf8=✓&authenticity_token=<token_value>&session[email]=username_here&session[password]=password_here&commit=Sign+In
Response Status: 200 OK
Response Header: Content-Type: text/html; charset=utf-8
Response Header: Set-Cookie: <cookie_value>
Response Body: <html>...</html>
```
