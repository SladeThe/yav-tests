package tests

import (
	"log"
	"reflect"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/SladeThe/yav"
	"github.com/SladeThe/yav/common"
	"github.com/go-playground/validator/v10"
)

var (
	playgroundValuelessCheckNames = map[string]struct{}{
		yav.CheckNameRequired:           {},
		yav.CheckNameRequiredWithAny:    {},
		yav.CheckNameRequiredWithoutAny: {},
		yav.CheckNameRequiredWithAll:    {},
		yav.CheckNameRequiredWithoutAll: {},
	}

	playgroundParameterlessCheckNames = map[string]struct{}{
		yav.CheckNameRequired: {},

		yav.CheckNameUnique: {},

		yav.CheckNameEmail: {},
		yav.CheckNameE164:  {},
		yav.CheckNameUUID:  {},

		yav.CheckNameLowercase: {},
		yav.CheckNameUppercase: {},

		yav.CheckNameContainsAlpha:      {},
		yav.CheckNameContainsLowerAlpha: {},
		yav.CheckNameContainsUpperAlpha: {},
		yav.CheckNameContainsDigit:      {},

		yav.CheckNameStartsWithAlpha:      {},
		yav.CheckNameStartsWithLowerAlpha: {},
		yav.CheckNameStartsWithUpperAlpha: {},
		yav.CheckNameStartsWithDigit:      {},

		yav.CheckNameEndsWithAlpha:      {},
		yav.CheckNameEndsWithLowerAlpha: {},
		yav.CheckNameEndsWithUpperAlpha: {},
		yav.CheckNameEndsWithDigit:      {},

		yav.CheckNameText:  {},
		yav.CheckNameTitle: {},
	}
)

type Playground struct {
	Validator *validator.Validate
}

func NewPlayground() Playground {
	p := Playground{Validator: validator.New()}

	// Register additional validation tags.
	p.mustRegisterContains()
	p.mustRegisterExcludes()
	p.mustRegisterStartsWith()
	p.mustRegisterEndsWith()
	p.mustRegisterText()
	p.mustRegisterTitle()

	return p
}

func (p Playground) Validate(s any) error {
	err := p.Validator.Struct(s)
	if err == nil {
		return nil
	}

	if fieldErrs, ok := err.(validator.ValidationErrors); ok {
		var yavErrs yav.Errors

		for _, fieldErr := range fieldErrs {
			tag := fieldErr.Tag()

			parameter := fieldErr.Param()
			if _, omitParameter := playgroundParameterlessCheckNames[tag]; omitParameter {
				parameter = ""
			}

			value := fieldErr.Value()
			if _, omitValue := playgroundValuelessCheckNames[tag]; omitValue {
				value = nil
			}

			yavErrs.Validation = append(yavErrs.Validation, yav.Error{
				CheckName: tag,
				Parameter: parameter,
				ValueName: fieldErr.Field(),
				Value:     value,
			})
		}

		return yavErrs.AsError()
	}

	return yav.Errors{Unknown: []error{err}}
}

func (p Playground) mustRegisterContains() {
	p.Validator.RegisterAlias(yav.CheckNameContainsAlpha, "containsany="+common.CharactersAlpha)
	p.Validator.RegisterAlias(yav.CheckNameContainsLowerAlpha, "containsany="+common.CharactersLowerAlpha)
	p.Validator.RegisterAlias(yav.CheckNameContainsUpperAlpha, "containsany="+common.CharactersUpperAlpha)
	p.Validator.RegisterAlias(yav.CheckNameContainsDigit, "containsany="+common.CharactersDigit)

	fn := func(fl validator.FieldLevel) bool {
		s, ok := p.fieldAsString(fl)
		if !ok || s == "" {
			return false
		}

		return strings.ContainsAny(s, common.CharactersSpecial)
	}

	if err := p.Validator.RegisterValidation(yav.CheckNameContainsSpecialCharacter, fn); err != nil {
		log.Fatal(err)
	}
}

func (p Playground) mustRegisterExcludes() {
	fn := func(fl validator.FieldLevel) bool {
		s, ok := p.fieldAsString(fl)
		if !ok {
			return false
		}

		for _, r := range s {
			if unicode.IsSpace(r) {
				return false
			}
		}

		return true
	}

	if err := p.Validator.RegisterValidation(yav.CheckNameExcludesWhitespace, fn); err != nil {
		log.Fatal(err)
	}
}

func (p Playground) mustRegisterStartsWith() {
	checks := []struct {
		name string
		do   func(r rune) bool
	}{{
		name: yav.CheckNameStartsWithAlpha,
		do:   common.IsAlpha,
	}, {
		name: yav.CheckNameStartsWithLowerAlpha,
		do:   common.IsLowerAlpha,
	}, {
		name: yav.CheckNameStartsWithUpperAlpha,
		do:   common.IsUpperAlpha,
	}, {
		name: yav.CheckNameStartsWithDigit,
		do:   common.IsDigit,
	}, {
		name: yav.CheckNameStartsWithSpecialCharacter,
		do:   common.IsSpecialCharacter,
	}}

	for _, check := range checks {
		do := check.do

		fn := func(fl validator.FieldLevel) bool {
			s, ok := p.fieldAsString(fl)
			if !ok || s == "" {
				return false
			}

			r, _ := utf8.DecodeRuneInString(s)
			return do(r)
		}

		if err := p.Validator.RegisterValidation(check.name, fn); err != nil {
			log.Fatal(err)
		}
	}
}

func (p Playground) mustRegisterEndsWith() {
	checks := []struct {
		name string
		do   func(r rune) bool
	}{{
		name: yav.CheckNameEndsWithAlpha,
		do:   common.IsAlpha,
	}, {
		name: yav.CheckNameEndsWithLowerAlpha,
		do:   common.IsLowerAlpha,
	}, {
		name: yav.CheckNameEndsWithUpperAlpha,
		do:   common.IsUpperAlpha,
	}, {
		name: yav.CheckNameEndsWithDigit,
		do:   common.IsDigit,
	}, {
		name: yav.CheckNameEndsWithSpecialCharacter,
		do:   common.IsSpecialCharacter,
	}}

	for _, check := range checks {
		do := check.do

		fn := func(fl validator.FieldLevel) bool {
			s, ok := p.fieldAsString(fl)
			if !ok || s == "" {
				return false
			}

			r, _ := utf8.DecodeLastRuneInString(s)
			return do(r)
		}

		if err := p.Validator.RegisterValidation(check.name, fn); err != nil {
			log.Fatal(err)
		}
	}
}

// mustRegisterText registers common.IsText as a string value validator.
func (p Playground) mustRegisterText() {
	fn := func(fl validator.FieldLevel) bool {
		s, ok := p.fieldAsString(fl)
		return ok && common.IsText(s)
	}

	if err := p.Validator.RegisterValidation(yav.CheckNameText, fn); err != nil {
		log.Fatal(err)
	}
}

// mustRegisterTitle registers common.IsTitle as a string value validator.
func (p Playground) mustRegisterTitle() {
	fn := func(fl validator.FieldLevel) bool {
		s, ok := p.fieldAsString(fl)
		return ok && common.IsTitle(s)
	}

	if err := p.Validator.RegisterValidation(yav.CheckNameTitle, fn); err != nil {
		log.Fatal(err)
	}
}

func (p Playground) RegisterRegexp(name, pattern string) error {
	rgx, errCompile := regexp.Compile(pattern)
	if errCompile != nil {
		return errCompile
	}

	fn := func(fl validator.FieldLevel) bool {
		s, ok := p.fieldAsString(fl)
		return ok && rgx.MatchString(s)
	}

	return p.Validator.RegisterValidation(name, fn)
}

func (p Playground) MustRegisterRegexp(name, pattern string) {
	if err := p.RegisterRegexp(name, pattern); err != nil {
		log.Fatal(err)
	}
}

func (p Playground) fieldAsString(fl validator.FieldLevel) (s string, ok bool) {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return "", false
	}

	return field.String(), true
}
