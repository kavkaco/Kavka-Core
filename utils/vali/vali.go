package vali

import (
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var vi = validator.New()

var locale_EN = en.New()
var uni = ut.New(locale_EN, locale_EN)
var ValiTranslator, _ = uni.GetTranslator("en")

var (
	valiLock     = &sync.Mutex{}
	valiInstance *Vali
)

type Vali struct{}

func Validator() *Vali {
	if valiInstance == nil {
		valiLock.Lock()
		defer valiLock.Unlock()

		valiInstance = &Vali{}
		valiInstance.translateOverride()
	}

	return valiInstance
}

func (*Vali) Validate(s interface{}) validator.ValidationErrors {
	err := vi.Struct(s)

	if err == nil {
		return validator.ValidationErrors{}
	}

	return err.(validator.ValidationErrors)
}

func (*Vali) translateOverride() {
	en_translations.RegisterDefaultTranslations(vi, ValiTranslator)

	vi.RegisterTranslation("required", ValiTranslator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must have a value by taha", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})
}
