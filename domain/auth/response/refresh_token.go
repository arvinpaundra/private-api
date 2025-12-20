package response

type RefreshToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserId       string `json:"user_id"`
}
