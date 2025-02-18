<div align="center">

<img src="https://raw.githubusercontent.com/0xStarLabs/TwitterAPI/public/X_logo.jpg" width="200" height="200">

# 🐦 Twitter API SDK

### A Modern, Feature-Rich Go SDK for Twitter API Automation

[![Stars](https://img.shields.io/github/stars/0xStarLabs/TwitterAPI?style=for-the-badge&logo=github&color=yellow)](https://github.com/0xStarLabs/TwitterAPI/stargazers)
![Twitter API](https://img.shields.io/badge/Twitter-API-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.23-00ADD8?style=for-the-badge&logo=go&logoColor=white)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)

<p align="center">
<b>Fast 🚀 Reliable 🛡️ Easy to Use 💡</b>
</p>

<p align="center">
A powerful Go SDK that provides seamless interaction with Twitter's API. Built with performance and ease of use in mind.
</p>

</div>

---

## Features

- 🔐 **Account Management**
  - Account validation
  - Support for auth tokens and JSON cookies
  - Proxy support
- 👥 **User Interactions**
  - Follow/Unfollow users
  - Get user information
  - Handle both usernames and user IDs
- 📝 **Content Creation**
  - Post tweets
  - Reply to tweets
  - Add media to posts
- ⚡ **Performance**
  - Efficient cookie management
  - Automatic CSRF token handling
  - Smart error handling

## Installation

```bash
go get github.com/0xStarLabs/TwitterAPI
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/0xStarLabs/TwitterAPI/client"
)

func main() {
    // Create a new account
    account := client.NewAccount("auth_token_here", "", "")
    
    // Initialize Twitter client
    twitter, err := client.NewTwitter(account)
    if err != nil {
        panic(err)
    }

    // Check if account is valid
    info, resp := twitter.IsValid()
    if resp.Success {
        fmt.Printf("Account %s is valid\n", info.Username)
    }
}
```

## Usage Examples

### Following Users

```go
// Follow by username
resp := twitter.Follow("username")
if resp.Success {
    fmt.Println("Successfully followed user")
}

// Unfollow by username or ID
resp = twitter.Unfollow("username")
if resp.Success {
    fmt.Println("Successfully unfollowed user")
}
```

### Posting Comments

```go
// Simple comment
resp := twitter.Comment("Great tweet!", "1234567890", nil)

// Comment with media
resp = twitter.Comment("Check this out!", "1234567890", &client.CommentOptions{
    MediaBase64: imageBase64,
})
```

### Getting User Info

```go
info, resp := twitter.GetUserInfoByUsername("username")
if resp.Success {
    fmt.Printf("User ID: %s\n", info.Data.User.Result.RestID)
    fmt.Printf("Followers: %d\n", info.Data.User.Result.Legacy.FollowersCount)
}
```

## Advanced Configuration

### Using Proxies

```go
account := client.NewAccount(
    "auth_token_here",
    "csrf_token",  // optional
    "user:pass@host:port",  // proxy
)
```

### Using JSON Cookies

```go
authToken, csrfToken, err := client.SetAuthCookies(
    0,  // account index
    cookieClient,
    `[{"name":"auth_token","value":"token"}]`,
)
```

## Error Handling

The SDK uses a consistent error handling pattern:

```go
resp := twitter.Follow("username")
if !resp.Success {
    switch resp.Status {
    case models.StatusAuthError:
        fmt.Println("Authentication failed")
    case models.StatusLocked:
        fmt.Println("Account is locked")
    default:
        fmt.Printf("Error: %v\n", resp.Error)
    }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This project is not affiliated with Twitter. Use at your own risk and ensure compliance with Twitter's Terms of Service.