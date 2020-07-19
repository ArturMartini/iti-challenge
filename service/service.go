package service

import (
	"github.com/arturmartini/iti-challenge/entities"
	log "github.com/sirupsen/logrus"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

type Service interface {
	ValidateStrongPassword(password entities.Password) bool
}

type service struct{}
type ValidateType func(string) bool

const (
	regexNotHaveSpace    = "^\\S*$"
	regexHaveOneDigit    = ".*\\d"
	regexMinCharacter    = ".*.{9,}"
	regexHaveLowerChar   = ".*[a-z]"
	regexHaveUpperChar   = ".*[A-Z]"
	regexHaveSpecialChar = ".*[!@#$%^&*()\\-+]"
)

var (
	once            sync.Once
	instance        Service
	regexValidation = []interface{}{
		regexp.MustCompile(regexNotHaveSpace),
		regexp.MustCompile(regexHaveOneDigit),
		regexp.MustCompile(regexMinCharacter),
		regexp.MustCompile(regexHaveLowerChar),
		regexp.MustCompile(regexHaveUpperChar),
		regexp.MustCompile(regexHaveSpecialChar),
		isNotRepeated,
	}
)

func New() Service {
	once.Do(func() {
		if instance == nil {
			instance = service{}
		}
	})
	return instance
}

func (r service) ValidateStrongPassword(password entities.Password) bool {
	valid := true
	for _, re := range regexValidation {
		if !executeValidation(re, password.Value) {
			valid = false
		}
	}
	return valid
}

func executeValidation(validate interface{}, password string) bool {
	valid := false
	if re, ok := validate.(*regexp.Regexp); ok {
		valid = re.MatchString(password)
	} else {
		if function, ok := validate.(func(string) bool); ok {
			valid = function(password)
		}
	}

	log.WithFields(log.Fields{
		"validate":      reflect.ValueOf(validate),
		"validate-type": reflect.TypeOf(validate),
		"valid":         valid}).Info("Execute validation")
	return valid
}

func isNotRepeated(value string) bool {
	chars := strings.Split(value, "")
	for ind, char := range chars {
		for repeatIndex, repeatChar := range chars {
			if char == repeatChar && ind != repeatIndex {
				return false
			}
		}
	}
	return true
}
