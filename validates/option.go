package validates

import (
	"database/sql/driver"
	"reflect"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// Option 验证器选项
type Option func(v *validator.Validate, trans ut.Translator)

// WithTag 设置Tag名称，默认：valid
func WithTag(s string) Option {
	return func(v *validator.Validate, trans ut.Translator) {
		v.SetTagName(s)
	}
}

// WithValuerType 注册自定义验证类型
func WithValuerType(types ...driver.Valuer) Option {
	customTypes := make([]any, 0, len(types))
	for _, t := range types {
		customTypes = append(customTypes, t)
	}

	return func(validate *validator.Validate, trans ut.Translator) {
		validate.RegisterCustomTypeFunc(func(field reflect.Value) any {
			if valuer, ok := field.Interface().(driver.Valuer); ok {
				val, _ := valuer.Value()
				return val
			}
			return nil
		}, customTypes...)
	}
}

// WithValidateFunc 注册自定义验证器
func WithValidateFunc(tag string, fn validator.Func, callValidationEvenIfNull ...bool) Option {
	return func(validate *validator.Validate, trans ut.Translator) {
		_ = validate.RegisterValidation(tag, fn, callValidationEvenIfNull...)
	}
}

// WithValidateFuncX 注册带Context的自定义验证器
func WithValidateFuncX(tag string, fn validator.FuncCtx, callValidationEvenIfNull ...bool) Option {
	return func(validate *validator.Validate, trans ut.Translator) {
		_ = validate.RegisterValidationCtx(tag, fn, callValidationEvenIfNull...)
	}
}

// WithTranslation 注册自定义错误翻译
//
//	参数 `text` 示例：
//	[示例1] {0}为必填字段
//	[示例2] {0}必须大于{1}
func WithTranslation(tag, text string, override bool) Option {
	return func(validate *validator.Validate, trans ut.Translator) {
		_ = validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
			return ut.Add(tag, text, override)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field(), fe.Param())
			return t
		})
	}
}
