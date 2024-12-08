package errors

type ErrOutPut struct{}

func (err ErrOutPut) Error() string {
	return "output error"
}

type ErrZeroSizeMatrix struct {
}

func (err ErrZeroSizeMatrix) Error() string {
	return "zero size matrix"
}
