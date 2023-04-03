package homework

type IntSliceValidator struct {
	name string
	val  any
	tag  string
}

type StrSliceValidator struct {
	name string
	val  any
	tag  string
}

func (v IntSliceValidator) Validate() ValidationErrors {
	value := v.val.([]int)
	errs := ValidationErrors{}

	for i := range value {
		e := IntValidator{v.name, value[i], v.tag}.Validate()
		errs = append(errs, e...)
	}

	return errs
}

func (v StrSliceValidator) Validate() ValidationErrors {
	value := v.val.([]string)
	errs := ValidationErrors{}

	for i := range value {
		e := StrValidator{v.name, value[i], v.tag}.Validate()
		errs = append(errs, e...)
	}

	return errs
}
