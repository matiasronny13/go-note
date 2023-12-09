package model

import "encoding/json"

type PageState struct {
	PageTitle        string  `json:"-"`
	PathId           string  `json:"-"`
	Id               string  `form:"id" json:"id,omitempty"`
	Content          string  `form:"content" json:"content"`
	Password         string  `form:"password" json:"password"`
	IsEncrypted      bool    `form:"is_encrypted" json:"is_encrypted"`
	IsEditMode       bool    `json:"-"`
	Errors           []error `json:"-"`
	ShowDeleteButton bool    `json:"-"`
}

func (r *PageState) String() []byte {
	data, _ := json.Marshal(r)
	return data
}
