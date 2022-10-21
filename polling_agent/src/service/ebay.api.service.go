package service

import (
	"encoding/json"
	"fmt"
	"github.com/jamespearly/loggly"
	"net/http"
	"net/url"
	"os"
	"polling_agent/handler"
	"polling_agent/model"
	"polling_agent/utils"
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
	TagRoot = "controller.ebay.controller.go."
)

func Search(q string, offset int, limit int, sort string) (response *http.Response, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", EbayRoot+F(EbaySearchEndpoint, q, offset, limit, sort), nil)
	req.Header.Set(EbayCMarketplaceIdKey, EbayCMarketplaceIdValueUs)
	req.Header.Set(ContentType, ApplicationJson)
	req.Header.Set(XEbayCEnduserctxKey, XEbayCEnduserctxValue)
	req.Header.Set(Authorization, os.Getenv("EBAY_BEARER_TOKEN"))
	response, err = client.Do(req)
	return
}

func TokenGenerator(wg *sync.WaitGroup) {

	// Instantiate Client
	tokenGeneratorTag := fmt.Sprintf("%stokenGenerator()", TagRoot)
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
	handler.HttpErrorCheck(err, tokenGenClient)

	// byte conversion and unmarshalling of token response info
	bodyBytes := utils.GetBytes(response, tokenGenClient)

	// initialization of the token model tokenModel
	tokenUnmarshalError := json.Unmarshal(bodyBytes, &tokenModel)
	handler.UnmarshalError(tokenUnmarshalError, tokenGenClient, string(bodyBytes))

	utils.SetTokenEnvs(tokenModel)
	wg.Done()
}

func RequestEbayBytes(ebayClient *loggly.ClientType, title string, sort string) []byte {

	response, httpErr := Search(title, 0, 3, sort)
	handler.HttpErrorCheck(httpErr, ebayClient)

	return utils.GetBytes(response, ebayClient)
}

func GenerateNewTokenIfNotExist() {
	var wg sync.WaitGroup

	if utils.TokenExpiredOrDoesntExist() {
		fmt.Println("authToken expired. generating new token. . .")
		wg.Add(1)
		TokenGenerator(&wg)
		wg.Wait()
		fmt.Println("authToken generated")
	}
}

func SearchItemsOnEbay(ebayClient *loggly.ClientType, malRes model.MalUserListResponse, sort string) []model.ItemSummaries {
	var ebayResponseModel model.EbaySearchResponse
	var itemSummaries []model.ItemSummaries
	for _, data := range malRes.Data {
		title := strings.ReplaceAll(
			utils.RemoveNonAlphaNums(data.Node.Title), " ", "+",
		)
		ebayBytes := RequestEbayBytes(ebayClient, title, sort)
		ebayUnmarshalErr := json.Unmarshal(ebayBytes, &ebayResponseModel)
		handler.UnmarshalError(ebayUnmarshalErr, ebayClient, title)
		itemSummaries = append(itemSummaries, ebayResponseModel.ItemSummaries...)
		time.Sleep(1 * time.Second)
	}
	return itemSummaries
}
