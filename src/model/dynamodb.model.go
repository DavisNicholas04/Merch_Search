package model

type UserEntry struct {
	UserId        string        `json:"user_id"`
	ItemId        string        `json:"item_id"`
	ItemSummaries ItemSummaries `json:"item_summaries"`
}
