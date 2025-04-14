package mux

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/itsoneiota/lambda-handlers/v2/pkg/helpers"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/serviceerror"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       string
}

func NewResponseWriter(w http.ResponseWriter, headers http.Header) *ResponseWriter {
	for k, v := range headers {
		if len(v) > 0 {
			w.Header().Set(k, v[0])
		}
	}

	return &ResponseWriter{
		ResponseWriter: w,
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *ResponseWriter) Write(body []byte) (int, error) {
	bodyStr := string(body)
	if !helpers.IsOkRange(w.statusCode) && !helpers.IsValidJSONObject(bodyStr) {
		var decodedString string
		if err := json.Unmarshal([]byte(bodyStr), &decodedString); err == nil {
			bodyStr = decodedString
		}

		e := serviceerror.NewServiceError(
			serviceerror.GetServiceErrorCode(w.statusCode),
			serviceerror.GetServiceErrorCode(w.statusCode),
			bodyStr,
		)

		b, err := json.Marshal(e)
		if err != nil {
			slog.Error(err.Error())
			return 0, err
		}

		bodyStr = string(b)
	}

	w.body = bodyStr
	w.ResponseWriter.Write([]byte(bodyStr))

	return len(body), nil
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *ResponseWriter) Body() string {
	return w.body
}

func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}
