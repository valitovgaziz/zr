package models

//go:generate reform

//reform:people
type News_db struct {
	Id      int64  `reform:"id,pk"`
	Title   string `reform:"title"`
	Content string `reform:"content"`
}
