package anigate

import "fmt"

func errString(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
