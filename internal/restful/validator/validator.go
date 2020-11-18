package validator

import "github.com/gin-gonic/gin/binding"

// ToV10 upgrade gin validator to v10
func ToV10() {
	binding.Validator = new(defaultValidator)
}
