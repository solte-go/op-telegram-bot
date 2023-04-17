package errhandler

import "fmt"

func Wrap(mgs string, err error) error {
	return fmt.Errorf("%s: %w", mgs, err)
}

func WrapIfErr(mgs string, err error) error {
	if err == nil {
		return nil
	}
	return Wrap(mgs, err)
}
