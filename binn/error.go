package binn

import "fmt"

const (
	CodeExpiredSubscription      = iota
	CodeNotFoundSubscription     = iota
	CodeNotFoundSubscribedBottle = iota

	CodeExpiredBottle     = iota
	CodeNotFoundBottle    = iota
	CodeUnavailableBottle = iota
)

type BinnError struct {
	Code int
	Msg  string
	err  error
}

func (err *BinnError) Error() string {
	if err.err != nil {
		return fmt.Sprintf("%s: %v", err.Msg, err.err)
	}
	return err.Msg
}

func NewBinnError(code int, msg string, err error) *BinnError {
	return &BinnError{
		Code: code,
		Msg:  msg,
		err:  err,
	}
}
