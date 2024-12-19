package commands

type CreateUserCommand struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DateOfBirth string `json:"date_of_birth"`
}
