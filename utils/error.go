package utils

import "fmt"

//CatchError catch panic
func CatchError(err *error) {
	if r := recover(); r != nil {
		switch r.(type) {
		case error:
			*err = r.(error)
		default:
			*err = fmt.Errorf("%v", r)
		}
	}
}
