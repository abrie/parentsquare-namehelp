# parentsquare-namehelp

## Login Method

This repository includes a method to log in to ParentSquare using an authenticity token, username, and password.

### Usage

1. Retrieve the authenticity token and cookie from the sessions page.
2. Call the `login` method with the authenticity token, username, and password.
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
