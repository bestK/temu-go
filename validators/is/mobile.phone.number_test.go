package is

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMobilePhoneNumber(t *testing.T) {
	tests := map[string]bool{
		"13401234567":    true,
		"13411234567":    true,
		"8613411234567":  true,
		" 8613411234567": true,
		"8713411234567":  false,
	}
	for str, ok := range tests {
		err := validation.Validate(str, validation.By(MobilePhoneNumber()))
		assert.Equal(t, ok, err == nil, str)
	}
}