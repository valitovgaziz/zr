package models

//go:generate reform

//reform:people
type News struct {
	Id         int64  `reform:"id,pk"`
	Title      string `reform:"title"`
	Content    string `reform:"content"`
	Categories []int  `reform:"categories"`
}
