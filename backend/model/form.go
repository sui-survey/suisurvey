package model

type CreateFormDto struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	ItemList []Item `json:"itemList"`
}

type GetFormDto struct {
	BlobId string `json:"blobId"`
}
