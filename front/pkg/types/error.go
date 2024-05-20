package types

import "fmt"

type STGError struct {
	Msg  string
	Code int
}

func (e *STGError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Msg)
}
