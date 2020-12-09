package validations

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

// Validate interface for things that want to be validatred
type Validate interface {
	Validate() error
}

func validate(db *gorm.DB) {
	if db.Statement.ReflectValue.CanAddr() {
		value := db.Statement.ReflectValue.Addr().Interface()
		_, err := govalidator.ValidateStruct(value)
		if err != nil {
			db.AddError(err)
		}
		if v, ok := value.(Validate); ok {
			if err := v.Validate(); err != nil {
				db.AddError(err)
			}
		}
	}
}

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *gorm.DB) {
	callback := db.Callback()
	if callback.Create().Get("validations:validate") == nil {
		callback.Create().Before("gorm:create").Register("validations:validate", validate)
	}
	if callback.Update().Get("validations:validate") == nil {
		callback.Update().Before("gorm:update").Register("validations:validate", validate)
	}
}
