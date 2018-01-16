package cassowary

type Variable struct {
	Value float64

	Name string

	owner *Param
}

func NewVariable(val float64) *Variable {
	return &Variable{
		Value: val,
	}
}

func (v *Variable) Param() *Param {
	return v.owner
}

func (v *Variable) applyUpdate(updated float64) bool {
	res := updated != v.Value
	v.Value = updated
	return res
}
