package entity

type UserAuth struct {
	UserId   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
}

func (u *UserAuth) GetUserId() string {
	return u.UserId
}

func (u *UserAuth) GetEmail() string {
	return u.Email
}

func (u *UserAuth) GetUsername() string {
	return u.Username
}

func (u *UserAuth) GetFullname() string {
	return u.Fullname
}
