package validate

//
// import (
// 	"errors"
// 	"reflect"
//
// 	"github.com/go-playground/locales/zh"
// 	ut "github.com/go-playground/universal-translator"
// 	"github.com/go-playground/validator/v10"
// 	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
// )
//
// var (
// 	uni      *ut.UniversalTranslator
// 	validate *validator.Validate
// 	trans    ut.Translator
// )
//
// func init() {
// 	local := zh.New()
// 	uni = ut.New(local, local)
// 	trans, _ = uni.GetTranslator("zh")
// 	validate = validator.New()
// 	_ = zhTranslations.RegisterDefaultTranslations(validate, trans)
// 	// 通过 json 标签获取字段名称
// 	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
// 		label := fld.Tag.Get("json")
// 		if label == "" {
// 			return fld.Name
// 		}
//
// 		return label
// 	})
// }
//
// func Validate(data any) error {
// 	if err := validate.Struct(data); err != nil {
//
// 		for _, err := range err.(validator.ValidationErrors) {
// 			return errors.New(err.Translate(trans))
// 		}
//
// 		return err
// 	}
//
// 	return nil
// }
