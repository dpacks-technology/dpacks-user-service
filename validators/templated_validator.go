package validators

import (
	"dpacks-go-services-template/models"
	"errors"
	"strings"
	"unicode"
)

func ValidateTemp(template models.TemplateModel, create bool) error {

	if !startsWithCapitalAndLimitedWords(template.Name, 5) {
		return errors.New("template Name should start with a capital letter and have no more than 5 words")
	}
	//Check if template description is not longer than 200 words and contains at least 1 sentence with more than 5 words
	if !isValidDescription(template.Description) {
		return errors.New("template Description should not be longer than 100 words and should contain at least 1 sentence with more than 5 words")
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
	if template.Price < 0 || template.Price > 200 {
		return errors.New("template Price should be between 0 and 100$")
	}

	return nil
}

func ValidateTempRating(templateRating models.TemplateRatingsModel) error {

	if templateRating.Rating < 1 || templateRating.Rating > 5 {
		return errors.New("template Rating should be between 1 and 5")
	}

	return nil
}

func ValidateTempEdit(template models.TemplateModel, create bool) error {

	if template.Name == "" {
		return errors.New("template Name cannot be empty")
	}
	if template.Category == "" {
		return errors.New("template Category cannot be empty")
	}
	if template.Description == "" {

		return errors.New("template Description cannot be empty")
	}
	if template.DevpDescription == "" {
		return errors.New("template Developer Message cannot be empty")
	}

	return nil
}

func startsWithCapitalAndLimitedWords(name string, maxWords int) bool {
	// Check if name is not empty
	if name == "" {
		return false
	}

	// Check if first character is uppercase
	if !unicode.IsUpper([]rune(name)[0]) {
		return false
	}

	// Split the name into words
	words := strings.Fields(name)

	// Check if number of words is within limit
	if len(words) > maxWords {
		return false
	}

	return true
}

func isValidDescription(description string) bool {
	// Check if description is not empty
	if description == "" {
		return false
	}

	// Check if description is not longer than 200 words
	if len(strings.Fields(description)) > 200 {
		return false
	}

	// Check if description contains at least 1 sentence with more than 5 words
	sentences := strings.Split(description, ".")
	for _, sentence := range sentences {
		words := strings.Fields(sentence)
		if len(words) > 5 {
			return true
		}
	}

	return false
}
