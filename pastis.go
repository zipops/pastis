package pastis

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
)

type Server struct{}

type HandlerFactory struct {
	encoder Encoder
	decoder Decoder
}

func (f HandlerFactory) WithEncoder(e Encoder) HandlerFactory {
	copyy := f
	copyy.encoder = e
	return copyy
}

func (f HandlerFactory) WithDecoder(d Decoder) HandlerFactory {
	copyy := f
	copyy.decoder = d
	return copyy
}

func implementsResponse(typ reflect.Type) bool {
	response := reflect.TypeOf((*Response)(nil)).Elem()
	return typ.Implements(response)
}

func (f HandlerFactory) Handler(handler interface{}) http.Handler {
	typ := reflect.TypeOf(handler)
	if typ.Kind() != reflect.Func {
		panic("handler must be a func")
	}

	reqIdx := -1
	httpResponseWriterIdx := -1
	httpRequestIdx := -1

	var reqTyp reflect.Type // must be a struct
	numIn := typ.NumIn()
	for i := 0; i < numIn; i++ {
		switch t := typ.In(i); t {
		case reflect.TypeOf((*http.ResponseWriter)(nil)).Elem():
			if httpResponseWriterIdx != -1 {
				panic("handler can take only one ResponseWriter parameter")
			}
			httpResponseWriterIdx = i
		case reflect.TypeOf((*http.Request)(nil)):
			if httpRequestIdx != -1 {
				panic("handler can take only one ResponseWriter parameter")
			}
			httpRequestIdx = i
		default:
			if t.Kind() == reflect.Struct {
				if reqIdx != -1 {
					panic("handler can take only one custom struct parameter")
				}
				reqIdx = i
				reqTyp = t
			} else {
				panic("handler can only take as parameter http.ResponseWriter, *http.Request or a struct")
			}
		}
	}

	reqField := make(map[string]int, 4)
	if reqIdx != -1 {
		reqTyp := typ.In(reqIdx)
		for i := 0; i < reqTyp.NumField(); i++ {
			f := reqTyp.Field(i)
			if f.Name == "Body" || f.Name == "URLParams" || f.Name == "URLQuery" || f.Name == "Header" {
				if f.Type.Kind() != reflect.Struct {
					panic(fmt.Sprint("field", f.Name, "must be a struct"))
				}
				reqField[f.Name] = i
			}
		}

	}

	numOut := typ.NumOut()
	if numOut > 1 || (numOut == 1 && !implementsResponse(typ.Out(0))) {
		panic("handler must return only one value implementing pastis.Response")
	}

	val := reflect.ValueOf(handler)
	return http.HandlerFunc(func(rw http.ResponseWriter, hreq *http.Request) {
		in := make([]reflect.Value, numIn)

		if reqIdx != -1 {
			req := reflect.New(reqTyp)

			if bodyIdx, ok := reqField["Body"]; ok {
				body := req.Elem().Field(bodyIdx).Addr().Interface()
				err := f.decoder.Decode(hreq.Body, body)
				if err != nil {
					sendError(f.encoder, rw, err)
					return
				}
			}
			in[reqIdx] = req.Elem()
		}

		if httpResponseWriterIdx != -1 {
			in[httpResponseWriterIdx] = reflect.ValueOf(rw)
		}

		if httpRequestIdx != -1 {
			in[httpRequestIdx] = reflect.ValueOf(hreq)
		}

		out := val.Call(in)

		if numOut == 1 {
			f.encoder.Encode(rw, out[0].Interface())
		}
	})
}

// Validator TODO
type Validator interface {
	Validate(interface{}) error
}

func sendError(encoder Encoder, rw http.ResponseWriter, err error) {
	var resp Response
	if perr, ok := err.(Error); ok {
		resp = perr
	} else {
		log.Println(err)
		resp = InternalError()
	}

	if err := sendResponse(encoder, rw, resp); err != nil {
		log.Println(err)
	}
}

func sendResponse(encoder Encoder, rw http.ResponseWriter, r Response) error {
	// if the user sets a headers it replaces the previous one
	if header := r.Header(); header != nil {
		for k, vs := range header {
			rw.Header().Del(k)
			for _, v := range vs {
				rw.Header().Add(k, v)
			}
		}
	}

	rw.WriteHeader(r.StatusCode())

	return encoder.Encode(rw, r)
}
