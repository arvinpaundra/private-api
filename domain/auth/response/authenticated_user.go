package response

type AuthenticatedUser struct {
	UserId   string `json:"user_id"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}
