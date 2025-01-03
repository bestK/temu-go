package is

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"strings"
)

// MobilePhoneOrTelNumber 判断是否为有效的手机或座机号码
func MobilePhoneOrTelNumber() validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("无效的手机/座机号码 %v", value)
		}

		if strings.TrimSpace(s) == "" {
			return errors.New("手机/座机号码为空")
		}

		if !mobilePhoneNumberPattern.MatchString(s) && !telNumberPattern.MatchString(s) {
			return fmt.Errorf("无效的手机/座机号码 %s", s)
		}

		return nil
	}
}
