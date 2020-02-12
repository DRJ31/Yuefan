package service

import (
	"fmt"
	"time"
)

type MyError struct {
	When time.Time
	Err  string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("Error: %v", e.Err)
}
