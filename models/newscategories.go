package models

//go:generate reform

//reform:people
type NewsCategories struct {
	CategoryId int64 `reform:"categoryid"`
	NewsId     int   `reform:"newsid"`
}
