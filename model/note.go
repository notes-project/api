package model

type Note struct {
	Title       string   `json:"title" binding:"required"`
	Date        string   `json:"date" binding:"-"`
	Description string   `json:"description" binding:"required"`
	Category    string   `json:"category" binding:"-"`
	Tags        []string `json:"tags" binding:"-"`
}
