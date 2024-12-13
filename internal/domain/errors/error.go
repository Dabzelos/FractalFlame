package errors

import "fmt"

type ErrOutPut struct {
	Err error
}

func (err ErrOutPut) Error() string {
	return fmt.Sprintf("output error occurred: %v", err.Err)
}

type ErrZeroSizeMatrix struct {
}

func (err ErrZeroSizeMatrix) Error() string {
	return "zero size matrix"
}

type ErrReadingConfig struct {
	Err error
}

func (err ErrReadingConfig) Error() string {
	return fmt.Sprintf("reading configuration error: %v", err.Err)
}

type ErrSavingImage struct {
	Err error
}

func (err ErrSavingImage) Error() string {
	return fmt.Sprintf("saving image error: %v", err.Err)
}
