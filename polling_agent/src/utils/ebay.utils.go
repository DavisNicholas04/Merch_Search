package utils

import (
	"fmt"
	"github.com/jamespearly/loggly"
	"merchSearch/handler"
	"merchSearch/model"
	"os"
	"strconv"
	"time"
)

const TagRoot = "controller.ebay.controller.go."

// SetTokenEnvs to be replaced with sending token info to a db
func SetTokenEnvs(tokenModel model.TokenInfo) {
	// Instantiate Client
	tokenGeneratorTag := fmt.Sprintf("%ssetTokenEnvs()", TagRoot)
	tokenGenClient := loggly.New(tokenGeneratorTag)

	setBearerErr := os.Setenv("EBAY_BEARER_TOKEN", fmt.Sprintf("Bearer %s", tokenModel.AccessToken))
	if setBearerErr != nil {
		clientErr := tokenGenClient.EchoSend("error", "Could not set Ebay Bearer token environment Variable")
		handler.ClientErrorCheck(clientErr)
	}
	// adds the time till expiration to the current time and converts it to a string in preparation of env/db storage
	dateOfTokenExpiration := strconv.FormatInt(
		time.Now().Add(time.Second*time.Duration(tokenModel.ExpiresIn)).Unix(),
		10,
	)
	setExpiryErr := os.Setenv("EBAY_BEARER_TOKEN_EXPIRATION", dateOfTokenExpiration)
	if setExpiryErr != nil {
		clientErr := tokenGenClient.EchoSend("error", "Could not set Ebay token timer token environment Variable")
		handler.ClientErrorCheck(clientErr)
	}
	setTokenTypeErr := os.Setenv("EBAY_TOKEN_TYPE", tokenModel.TokenType)
	if setTokenTypeErr != nil {
		clientErr := tokenGenClient.EchoSend("error", "Could not set Ebay token type environment Variable")
		handler.ClientErrorCheck(clientErr)
	}
}

func TokenExpiredOrDoesntExist() bool {
	expirationTime, _ := strconv.ParseInt(os.Getenv("EBAY_BEARER_TOKEN_EXPIRATION"), 10, 64)
	EbayTokenExpired := expirationTime < time.Now().Unix()
	TokenNotExist := !(Second(os.LookupEnv("EBAY_BEARER_TOKEN")) &&
		Second(os.LookupEnv("EBAY_BEARER_TOKEN_EXPIRATION")) &&
		Second(os.LookupEnv("EBAY_TOKEN_TYPE")))
	return EbayTokenExpired || TokenNotExist
}
