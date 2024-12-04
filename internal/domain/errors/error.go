package errors

type ErrOutPut struct{}

func (err ErrOutPut) Error() string {
	return "output error"
}
