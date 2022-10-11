package controller

import (
	"encoding/json"
	"fmt"
	"github.com/jamespearly/loggly"
	"merchSearch/model"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// Common HEADER KEYS

	Authorization = "Authorization"
	ContentType   = "Content-Type"

	// Common Header Values

	ApplicationJson    = "application/json"
	XWwwFormUrlencoded = "application/x-www-form-urlencoded"

	// search function constant variables

	EbayRoot           = "https://api.sandbox.ebay.com/buy/browse/v1/"
	EbaySearchEndpoint = "item_summary/search?q=%s+anime&offset=%d&limit=%d&sort=%s" // EbaySearchEndpoint URI parameters: q, offset, limit, sort

	// Ebay Search Headers

	EbayCMarketplaceIdKey     = "X-EBAY-C-MARKETPLACE-ID"
	EbayCMarketplaceIdValueUs = "EBAY_US"
	XEbayCEnduserctxKey       = "X-EBAY-C-ENDUSERCTX"
	XEbayCEnduserctxValue     = "affiliateCampaignId=<ePNCampaignId>,affiliateReferenceId=<referenceId>"

	// tokenGenerator function constant variables

	EbayTokenGeneratorEndpoint = "https://api.sandbox.ebay.com/identity/v1/oauth2/token?"

	// Ebay Oauth Token Generator Body

	GrantType         = "grant_type"
	ClientCredentials = "client_credentials"
	Scope             = "scope"
	ScopeValue        = "https://api.ebay.com/oauth/api_scope"

	// tagRoot
	tagRoot = "controller.ebay.controller.go."
)

func search(q string, offset int, limit int, sort string) (response *http.Response, err error) {
	var f = fmt.Sprintf

	client := &http.Client{}
	req, _ := http.NewRequest("GET", EbayRoot+f(EbaySearchEndpoint, q, offset, limit, sort), nil)
	req.Header.Set(EbayCMarketplaceIdKey, EbayCMarketplaceIdValueUs)
	req.Header.Set(ContentType, ApplicationJson)
	req.Header.Set(XEbayCEnduserctxKey, XEbayCEnduserctxValue)
	req.Header.Set(Authorization, os.Getenv("EBAY_BEARER_TOKEN"))
	response, err = client.Do(req)
	return
}

func tokenGenerator(wg *sync.WaitGroup) {

	// Instantiate Client
	tokenGeneratorTag := fmt.Sprintf("%stokenGenerator()", tagRoot)
	tokenGenClient := loggly.New(tokenGeneratorTag)

	// declaration of a token model
	var tokenModel model.TokenInfo

	// set tokenGenerator body
	unencodedBody := url.Values{
		GrantType: {ClientCredentials},
		Scope:     {ScopeValue},
	}
	encodedBody := unencodedBody.Encode()

	// make http request to ebay oauth generator endpoint
	client := &http.Client{}
	req, _ := http.NewRequest("POST", EbayTokenGeneratorEndpoint, strings.NewReader(encodedBody))
	req.Header.Set(ContentType, XWwwFormUrlencoded)
	req.Header.Set(Authorization, os.Getenv("EBAY_CREDENTIALS"))
	response, err := client.Do(req)
	httpErrorCheck(err, tokenGenClient)

	// byte conversion and unmarshalling of token response info
	bodyBytes := GetBytes(response, tokenGenClient)

	// initialization of the token model tokenModel
	tokenUnmarshalError := json.Unmarshal(bodyBytes, &tokenModel)
	unmarshalError(tokenUnmarshalError, tokenGenClient)

	setTokenEnvs(tokenModel)
	wg.Done()
}

// to be replaced with sending token info to a db
func setTokenEnvs(tokenModel model.TokenInfo) {
	// Instantiate Client
	tokenGeneratorTag := fmt.Sprintf("%ssetTokenEnvs()", tagRoot)
	tokenGenClient := loggly.New(tokenGeneratorTag)

	setBearerErr := os.Setenv("EBAY_BEARER_TOKEN", fmt.Sprintf("Bearer %s", tokenModel.AccessToken))
	if setBearerErr != nil {
		clientErr := tokenGenClient.EchoSend("error", "Could not set Ebay Bearer token environment Variable")
		clientErrorCheck(clientErr)
	}
	// adds the time till expiration to the current time and converts it to a string in preparation of env/db storage
	dateOfTokenExpiration := strconv.FormatInt(
		time.Now().Add(time.Second*time.Duration(tokenModel.ExpiresIn)).Unix(),
		10,
	)
	setExpiryErr := os.Setenv("EBAY_BEARER_TOKEN_EXPIRATION", dateOfTokenExpiration)
	if setExpiryErr != nil {
		clientErr := tokenGenClient.EchoSend("error", "Could not set Ebay token timer token environment Variable")
		clientErrorCheck(clientErr)
	}
	setTokenTypeErr := os.Setenv("EBAY_TOKEN_TYPE", tokenModel.TokenType)
	if setTokenTypeErr != nil {
		clientErr := tokenGenClient.EchoSend("error", "Could not set Ebay token type environment Variable")
		clientErrorCheck(clientErr)
	}
}
