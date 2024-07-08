package vali

import (
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var vi = validator.New()

var (
	locale_EN         = en.New()
	uni               = ut.New(locale_EN, locale_EN)
	ValiTranslator, _ = uni.GetTranslator("en")
)

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
	en_translations.RegisterDefaultTranslations(vi, ValiTranslator) // nolint
}
