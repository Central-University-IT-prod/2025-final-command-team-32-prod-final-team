package service

import (
	"reflect"
	"regexp"

	"solution/internal/domain/contracts"
	"solution/internal/domain/dto"
	"solution/internal/wrapper"

	"github.com/go-playground/validator/v10"
)

var _ contracts.ValidatorService = (*validatorService)(nil)

type validatorService struct {
	validator *validator.Validate
}

func NewValidatorService() *validatorService {
	service := validatorService{}
	service.validator = validator.New(validator.WithRequiredStructEnabled())
	service.initTags()
	return &service
}

func (s *validatorService) initTags() {
	s.validator.RegisterValidation("username", usernameTag)
	s.validator.RegisterValidation("password", passwordTag)
}

// dto must be a pointer
func (s *validatorService) ValidateRequestData(dto any) *dto.HttpErr {
	err := s.validator.Struct(dto)
	if err != nil {
		return wrapper.ValidationErr(err.Error())
	}
	return nil
}

func (s *validatorService) ValidateRequestSlice(dtos any) *dto.HttpErr {
	// no additional validation required (type assertion of dataModel being a slice is already executed in appropriate handler with fiber.BodyParser())
	// usage without dataModels type assertion being a slice is prohibited!
	v := reflect.ValueOf(dtos)
	for i := range v.Len() {
		err := s.validator.Struct(v.Index(i).Interface())
		if err != nil {
			return wrapper.ValidationErr(err.Error())
		}
	}
	return nil
}

func usernameTag(fl validator.FieldLevel) bool {
	username := fl.Field()
	if username.Len() < 3 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-zA-Zа-яА-Я0-9.,!?:;()\s"'«»—–-]+$`, username.String())
	return matched
}

func passwordTag(fl validator.FieldLevel) bool {
	password := fl.Field()
	return password.Len() >= 8
}
