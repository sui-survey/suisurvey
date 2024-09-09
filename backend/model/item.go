package model

import "kevinsheeran/walrus-backend/constant/enums"

type Item struct {
	Title string          `json:"title"`
	Name  string          `json:"name"`
	Type  enums.ItermType `json:"type"`
	Value string          `json:"value"`
}
