package validate

import (
	dnjinminUt "gitee.com/DNjinmin/universal-translator"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
	zh_translation "github.com/go-playground/validator/v10/translations/zh"
)

var Trans ut.Translator

type Validate interface {
	Tag() string
	ValidateFn() validator.Func
	Error() string
}

func RegisterValidatorFunc(v *validator.Validate, tag string, msgStr string, fn func(fl validator.FieldLevel) bool) {
	// 先注册验证器
	v.RegisterValidation(tag, fn)
	// 自定义错误的内容
	v.RegisterTranslation(tag, Trans, func(ut ut.Translator) error {
		return ut.Add(tag, msgStr, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}

func InitValidate(v *validator.Validate, validators []Validate, locale string) {
	if v == nil {
		return
	}

	// 注册一个对咋们错误信息翻译的工具
	if ok := translatorZh(locale, v); !ok {
		return
	}
	for _, val := range validators {
		RegisterValidatorFunc(v, val.Tag(), val.Error(), val.ValidateFn())
	}
}

func translatorZh(locale string, v *validator.Validate) (ok bool) {
	// 创建出中文和英文翻译器
	zhT := zh.New()
	enT := en.New()
	// 构建一个语言环境; 第一个参数是语言环境，第二个翻译的语言，最后是应该支持的环境
	uni := dnjinminUt.New(enT, zhT, enT)
	Trans, ok = uni.GetTranslator(locale)
	if !ok {
		return ok
	}
	switch locale {
	case "en":
		en_translation.RegisterDefaultTranslations(v, Trans)
	case "zh":
		zh_translation.RegisterDefaultTranslations(v, Trans)
	default:
	}

	return true
}
