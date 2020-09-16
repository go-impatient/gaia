package service

type Services struct {
	User  UserService
	Admin AdminService
}

func NewServices(User UserService, Admin AdminService) *Services {
	return &Services{
		User,
		Admin,
	}
}
