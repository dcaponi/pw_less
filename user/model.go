package user

type User struct {
	ID    int64  `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}
