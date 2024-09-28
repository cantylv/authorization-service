package user

import (
	"github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/entity/dto"
)

func newUserFromSignUpForm(data *dto.CreateData, hashedPassword string) *entity.User {
	return &entity.User{
		Email:     data.Email,
		Password:  hashedPassword,
		FirstName: data.FirstName,
		LastName:  data.LastName,
	}
}
