package dto

import (
	"regexp"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	me "github.com/cantylv/authorization-service/microservices/task_manager/internal/utils/myerrors"
)

var (
	nameRegexp     = regexp.MustCompile(`^[A-ZА-ЯЁ][a-zA-Zа-яА-ЯёЁ\s-]{1,50}$`)
	passwordRegexp = regexp.MustCompile(`^[A-Za-z\d!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]+$`)
)

// INPUT DATAFLOW
type CreateData struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (h *CreateData) Validate() error {
	if !govalidator.IsEmail(h.Email) {
		return me.ErrInvalidEmail
	}
	if isMatch := nameRegexp.MatchString(h.FirstName); !isMatch {
		return me.ErrInvalidFirstName
	}
	if isMatch := nameRegexp.MatchString(h.LastName); !isMatch {
		return me.ErrInvalidLastName
	}
	return isPasswordValid(h.Password)
}

func isPasswordValid(pwd string) error {
	pwdLen := utf8.RuneCountInString(pwd)
	if pwdLen > 30 {
		return me.ErrPasswordTooLong
	}
	if pwdLen < 8 {
		return me.ErrPasswordTooShort
	}

	// проверяем наличие хотя бы одной буквы верхнего регистра
	letterRegex := regexp.MustCompile(`[A-Z]`)
	if !letterRegex.MatchString(pwd) {
		return me.ErrPasswordFormat
	}

	// проверяем наличие хотя бы одной цифры
	digitRegex := regexp.MustCompile(`\d`)
	if !digitRegex.MatchString(pwd) {
		return me.ErrPasswordFormat
	}

	if !passwordRegexp.MatchString(pwd) {
		return me.ErrPasswordFormat
	}
	return nil
}

// OUTPUT DATAFLOW
type UserWithoutPassword struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
