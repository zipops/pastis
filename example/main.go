package main

import (
	"log"
	"net/http"

	"github.com/zipops/pastis"
)

type azeRequest struct {
	Body struct {
		Toto string `json:"toto"`
	}
}

type azeResponse struct {
	pastis.GenericResponse
	Titi string `json:"titi"`
}

func azeHandler(r azeRequest) azeResponse {
	return azeResponse{
		Titi: r.Body.Toto,
	}
}

func main() {
	factory := (pastis.HandlerFactory{}).
		WithDecoder(pastis.EncodingJSON{}).
		WithEncoder(pastis.EncodingJSON{})

	log.Fatal(http.ListenAndServe(":8080", factory.Handler(azeHandler)))
	// http POST http://localhost:8080/ toto=salut
}
