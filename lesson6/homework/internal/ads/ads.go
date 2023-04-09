package ads

type Ad struct {
	ID        int64
	Title     string `validate:"min:1 max:100"`
	Text      string `validate:"min:1 max:500"`
	AuthorID  int64
	Published bool
}
