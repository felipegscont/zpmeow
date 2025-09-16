package validation

type Validator = BaseValidator

func NewValidator() *Validator {
	return NewBaseValidator()
}
