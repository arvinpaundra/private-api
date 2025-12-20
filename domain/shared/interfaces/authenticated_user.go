package interfaces

type AuthenticatedUser interface {
	GetUserId() string
	GetUsername() string
	GetEmail() string
	GetFullname() string
}
