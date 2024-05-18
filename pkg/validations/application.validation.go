package validations

type ApplicationRegisterBody struct {
	Name        string `json:"name" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Description string `json:"description" validate:"required"`
}
