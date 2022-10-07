package controller

import (
	"fmt"
	"net/http"
	"os"
)

const (
	EbayRoot                  = "https://api.sandbox.ebay.com/buy/browse/v1/"
	EbayCMarketplaceIdKey     = "X-EBAY-C-MARKETPLACE-ID"
	EbayCMarketplaceIdValueUs = "EBAY_US"
	ContentTypeKey            = "Content-Type"
	ContentTypeValue          = "application/json"
	xEbayCEnduserctxKey       = "X-EBAY-C-ENDUSERCTX"
	xEbayCEnduserctxValue     = "affiliateCampaignId=<ePNCampaignId>,affiliateReferenceId=<referenceId>"

	EbaySearchEndpoint = "item_summary/search?q=%s+anime&offset=%d&limit=%d&sort=%s" // EbaySearchEndpoint URI parameters: q, offset, limit, sort
)

var (
	f = fmt.Sprintf
)

func search(q string, offset int, limit int, sort string) (response *http.Response, err error) {
	client := &http.Client{}
	fmt.Sprintf("%d", 3)
	req, _ := http.NewRequest("GET", EbayRoot+f(EbaySearchEndpoint, q, offset, limit, sort), nil)
	req.Header.Set(EbayCMarketplaceIdKey, EbayCMarketplaceIdValueUs)
	req.Header.Set(ContentTypeKey, ContentTypeValue)
	req.Header.Set(xEbayCEnduserctxKey, xEbayCEnduserctxValue)
	req.Header.Set("Authorization", os.Getenv("EBAY_BEARER_TOKEN"))
	response, err = client.Do(req)
	return
}
