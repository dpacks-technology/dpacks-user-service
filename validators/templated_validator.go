package validators

import (
	"dpacks-go-services-template/models"
	"errors"
)

func ValidateTemp(template models.TemplateModel, create bool) error {

	if template.Name == "" {
		return errors.New("template Name cannot be empty")
	}
	if template.Description == "" {
		return errors.New("template Description cannot be empty")
	}
	if template.Category == "" {
		return errors.New("template Category cannot be empty")
	}
	if template.MainFile == "" {
		return errors.New("template MainFile cannot be empty")
	}
	if template.ThmbnlFile == "" {
		return errors.New("template ThmbnlFile cannot be empty")
	}

	return nil
}
