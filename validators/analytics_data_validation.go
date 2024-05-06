package validators

import "github.com/go-playground/validator/v10"

//validate phonenumber is correct or not

func NumberValidation(Number int) bool {
	v := validator.New()

	err := v.Var(Number, "required,numeric")
	if err != nil {
		return false
	} else {
		return true
	}

}

// string validation can only contain strings
func StringValidation(str string) bool {
	v := validator.New()

	err := v.Var(str, "required,alpha")
	if err != nil {
		return false
	} else {
		return true
	}
}
