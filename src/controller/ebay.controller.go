package controller

import (
	"fmt"
	"net/http"
)

const (
	EbayRoot                  = "https://api.sandbox.ebay.com/buy/browse/v1/"
	EbayCMarketplaceIdKey     = "X-EBAY-C-MARKETPLACE-ID"
	EbayCMarketplaceIdValueUs = "EBAY_US"
	ContentTypeKey            = "Content-Type"
	ContentTypeValue          = "application/json"
	xEbayCEnduserctxKey       = "X-EBAY-C-ENDUSERCTX"
	xEbayCEnduserctxValue     = "affiliateCampaignId=<ePNCampaignId>,affiliateReferenceId=<referenceId>"

	EbaySearchEndpoint = "item_summary/search?&%s%s%s%s" // EbaySearchEndpoint URI parameters: q, offset, limit, sort
)

var (
	f = fmt.Sprintf
)

func search(q string, offset string, limit string, sort string) (response *http.Response, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", EbayRoot+f(EbaySearchEndpoint, q, offset, limit, sort), nil)
	req.Header.Set(EbayCMarketplaceIdKey, EbayCMarketplaceIdValueUs)
	req.Header.Set(ContentTypeKey, ContentTypeValue)
	req.Header.Set(xEbayCEnduserctxKey, xEbayCEnduserctxValue)
	response, err = client.Do(req)
	return
}
