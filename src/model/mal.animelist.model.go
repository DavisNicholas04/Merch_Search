package model

type MalUserListResponse struct {
	Data   []Data `json:"data"`
	Paging Paging `json:"paging"`
}

type Data struct {
	Node Node `json:"node"`
}

type Node struct {
	Id          int32       `json:"id"`
	Title       string      `json:"title"`
	MainPicture MainPicture `json:"main_picture"`
}

type MainPicture struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type Paging struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}
