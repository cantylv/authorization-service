package user

import (
	ent "github.com/cantylv/authorization-service/internal/entity"
	"github.com/cantylv/authorization-service/internal/entity/dto"
)

func getUserWithoutPassword(user *ent.User) *dto.UserWithoutPassword {
	return &dto.UserWithoutPassword{
		ID:        user.Id,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}
