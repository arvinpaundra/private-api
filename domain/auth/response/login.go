package response

type UserLogin struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserId       string `json:"user_id"`
}
