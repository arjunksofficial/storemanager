package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"storemanager/internal/common/responsehelper"
)

// RequestBodyParser validates the user to access handler
func RequestBodyParser(v interface{}) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			contentType := r.Header.Get("Content-Type")
			var decoder Decoder
			a := reflect.New(reflect.TypeOf(v)).Interface()
			if contentType == "" {
				responsehelper.RespondWithErrorAsJSON(400, w, "content_type should not be empty")
				return
			}
			switch contentType {
			case "application/json":
				decoder = JSONType{}
			default:
				responsehelper.RespondWithErrorAsJSON(400, w, "content_type not supported")
				return
			}
			if r.Body == nil {
				responsehelper.RespondWithErrorAsJSON(400, w, "nil body")
				return
			}
			err := decoder.Decode(r.Body, a)
			if err != nil {
				responsehelper.RespondWithErrorAsJSON(400, w, err.Error())
				return
			}
			ctx = context.WithValue(ctx, ParsedRequest, a)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}

}

// JSONType is struct used to parse JSON request bosy
type JSONType struct{}

// Decoder interface for decoding http requsest body according to content-type
type Decoder interface {
	Decode(io.ReadCloser, interface{}) error
}

// Decode decodes request body and gives a parsed struct object
func (d JSONType) Decode(reqBody io.ReadCloser, out interface{}) (err error) {
	err = json.NewDecoder(reqBody).Decode(out)
	return err
}
