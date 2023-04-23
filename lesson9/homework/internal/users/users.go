package users

type User struct {
	ID       int64
	Nickname string `validate:"min:1 max:100"`
	Email    string `validate:"min:0 max:500"`
}
