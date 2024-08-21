package validator

import (
	"reflect"

	"github.com/go-playground/locales/zh_Hans_CN"
	unTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
)

func Validate(data interface{}, language ...string) error {
	validate := validator.New()

	var trans unTrans.Translator
	if len(language) > 0 {
		switch language[0] {
		case "zh_Hans_CN":
			uni := unTrans.New(zh_Hans_CN.New())
			trans, _ = uni.GetTranslator(language[0])
			err := zhTrans.RegisterDefaultTranslations(validate, trans)
			if err != nil {
				return err
			}
		}
	}

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			label = field.Tag.Get("json")
			if label == "" {
				label = field.Tag.Get("form")
				if label == "" {
					label = field.Tag.Get("path")
				}
			}
		}
		return label
	})

	err := validate.Struct(data)
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			if trans != nil {
				return errors.Errorf(v.Translate(trans))
			}
			return v
		}
	}
	return nil
}
