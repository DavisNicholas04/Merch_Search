package model

type EbaySearchResponse struct {
	Href          string          `json:"href"`
	Total         int             `json:"total"`
	Next          string          `json:"next"`
	Prev          string          `json:"prev"`
	Limit         int             `json:"limit"`
	Offset        int             `json:"offset"`
	ItemSummaries []ItemSummaries `json:"itemSummaries"`
}

type Categories struct {
	CategoryId   string `json:"categoryId"`
	CategoryName string `json:"categoryName"`
}

type Image struct {
	ImageURL string `json:"imageUrl"`
}

type Price struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Seller struct {
	Username           string `json:"username"`
	FeedbackPercentage string `json:"feedbackPercentage"`
	FeedbackScore      int    `json:"feedbackScore"`
}

type ThumbnailImages struct {
	ImageURL string `json:"imageUrl"`
}

type ShippingCost struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type ShippingOptions struct {
	ShippingCostType string       `json:"shippingCostType"`
	ShippingCost     ShippingCost `json:"shippingCost"`
}

type ItemLocation struct {
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type AdditionalImages struct {
	ImageURL string `json:"imageUrl"`
}

type ItemSummaries struct {
	ItemId           string             `json:"itemId"`
	Title            string             `json:"title"`
	Categories       []Categories       `json:"categories"`
	Image            Image              `json:"image"`
	Price            Price              `json:"price"`
	Seller           Seller             `json:"seller"`
	Condition        string             `json:"condition"`
	ThumbnailImages  []ThumbnailImages  `json:"thumbnailImages"`
	ItemWebURL       string             `json:"itemWebUrl"`
	AdditionalImages []AdditionalImages `json:"additionalImages"`
	AdultOnly        bool               `json:"adultOnly"`
}
