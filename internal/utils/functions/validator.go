package functions

import "github.com/asaskevich/govalidator"

func Validate(data any) (bool, error) {
	return govalidator.ValidateStruct(data)
}
