package entity

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Session struct {
	UserID int
	Token  string
}
