package validators

import (
	"fmt"

	"github.com/markbates/going/validate"
)

type IntIsPresent struct {
	Name  string
	Field int
}

func (v *IntIsPresent) IsValid(errors *validate.Errors) {
	if v.Field == 0 {
		errors.Add(generateKey(v.Name), fmt.Sprintf("%s can not be blank.", v.Name))
	}
}
