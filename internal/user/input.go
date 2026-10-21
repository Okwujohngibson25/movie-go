package user

type CreateUserInput struct {
	Fullname string
	Email    string
	Password string
}

type LoginUserInput struct {
	Email    string
	Password string
}
