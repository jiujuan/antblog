// Package validator 注册业务自定义校验规则。
package validator

import (
	"regexp"

	gvalidator "github.com/go-playground/validator/v10"

	"antblog/pkg/validator"
)

var mobileRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// RegisterAll 注册所有自定义校验规则
func RegisterAll() error {
	v := validator.Default()

	// 手机号校验
	if err := v.RegisterRule("mobile", validateMobile, "{0}必须是有效的手机号"); err != nil {
		return err
	}

	// 用户名格式：字母/数字/下划线/连字符，3-32 位
	if err := v.RegisterRule("username", validateUsername, "{0}只能包含字母、数字、下划线和连字符，长度3-32位"); err != nil {
		return err
	}

	// 安全密码：至少包含字母和数字
	if err := v.RegisterRule("safe_password", validateSafePassword, "{0}必须包含字母和数字"); err != nil {
		return err
	}

	return nil
}

func validateMobile(fl gvalidator.FieldLevel) bool {
	return mobileRegex.MatchString(fl.Field().String())
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]{3,32}$`)

func validateUsername(fl gvalidator.FieldLevel) bool {
	return usernameRegex.MatchString(fl.Field().String())
}

func validateSafePassword(fl gvalidator.FieldLevel) bool {
	pwd := fl.Field().String()
	hasLetter := false
	hasDigit := false
	for _, c := range pwd {
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z':
			hasLetter = true
		case c >= '0' && c <= '9':
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}
