package pastis

import (
	"fmt"
	"net/http"
)

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.String()
}

func (e Error) String() string {
	return fmt.Sprintf("pastis http error: [%d %s] %s", e.Status, http.StatusText(e.Status), e.Message)
}

func (e Error) Header() http.Header {
	return nil
}

func (e Error) StatusCode() int {
	return e.Status
}

func InternalError() Error {
	return Error{
		Status:  500,
		Message: "Internal Server Error",
	}
}

func Err(status int, msg string) Error {
	return Error{
		Status:  status,
		Message: msg,
	}
}
