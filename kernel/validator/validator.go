package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslation "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

var (
	uni   *ut.UniversalTranslator
	valid *validator.Validate
	trans ut.Translator
)

func Init() {

	//注册翻译器
	chinese := zh.New()
	uni = ut.New(chinese, chinese)

	trans, _ = uni.GetTranslator("zh")

	//获取gin的校验器
	valid = binding.Validator.Engine().(*validator.Validate)

	valid.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("label")
	})

	_ = valid.RegisterValidation("mobile", mobile)
	_ = valid.RegisterValidation("dir", dir)
	_ = valid.RegisterValidation("username", username)
	_ = valid.RegisterValidation("password", password)
	_ = valid.RegisterValidation("snowflake", snowflake)

	_ = valid.RegisterTranslation("mobile", trans, func(ut ut.Translator) error {
		return ut.Add("mobile", "手机号格式错误", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mobile")
		return t
	})

	_ = valid.RegisterTranslation("dir", trans, func(ut ut.Translator) error {
		return ut.Add("dir", "文件夹格式错误", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("dir")
		return t
	})

	_ = valid.RegisterTranslation("username", trans, func(ut ut.Translator) error {
		return ut.Add("username", "请输入 4-20 位的英文字母数字以及 -_ 等字符", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("username")
		return t
	})

	_ = valid.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("password", "请输入 6-32 位的英文字母数字以及 -_@$&%! 等特殊字符", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password")
		return t
	})

	_ = valid.RegisterTranslation("snowflake", trans, func(ut ut.Translator) error {
		return ut.Add("snowflake", "雪花 ID 格式错误", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("snowflake")
		return t
	})

	//注册翻译器
	_ = zhTranslation.RegisterDefaultTranslations(valid, trans)
}

//Translate 翻译错误信息
func Translate(err error) string {

	result := ""

	errors := err.(validator.ValidationErrors)

	for _, item := range errors {
		result = item.Translate(trans)
		break
	}

	return result
}

func Translates(err error) map[string][]string {

	var result = make(map[string][]string)

	errors := err.(validator.ValidationErrors)

	for _, item := range errors {
		result[item.Field()] = append(result[item.Field()], item.Translate(trans))
	}

	return result
}
