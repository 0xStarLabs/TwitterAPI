package main

import (
	"log"
	"fmt"
	
	"github.com/0xStarLabs/TwitterAPI/client"
	"github.com/0xStarLabs/TwitterAPI/models"
	"github.com/0xStarLabs/TwitterAPI/utils"
)

func main() {
	proxy := "user:pass@host:port"
	authToken := "auth_token_here"
	// Create account
	account := client.NewAccount(authToken, "", proxy)

	// Or with detailed logging
	verboseConfig := models.NewConfig()
	verboseConfig.LogLevel = utils.LogLevelDebug

	twitter, err := client.NewTwitter(account, verboseConfig)
	if err != nil {
		// Handle error your way
		log.Fatal(err)
	}

	info, status := twitter.IsValid()
	if status.Error != nil {
		log.Fatal(status.Error)
	}
	fmt.Println(info)
}
