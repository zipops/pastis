package pastis

import "net/http"

type Response interface {
	Header() http.Header
	StatusCode() int
}

type GenericResponse struct {
	header http.Header
	status int
}

func (r GenericResponse) SetHeader(h http.Header) {
	r.header = h
}

func (r GenericResponse) Header() http.Header {
	return r.header
}

func (r GenericResponse) StatusCode() int {
	return r.status
}
