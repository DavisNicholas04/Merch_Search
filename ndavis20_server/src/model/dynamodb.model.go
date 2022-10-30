package model

type UserEntry struct {
	UserId        string        `json:"user_id"`
	ItemId        string        `json:"item_id"`
	ItemSummaries ItemSummaries `json:"item_summaries"`
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

type AdditionalImages struct {
	ImageURL string `json:"imageUrl"`
}
