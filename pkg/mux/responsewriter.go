package mux

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/itsoneiota/lambda-handlers/v2/pkg/helpers"
	"github.com/itsoneiota/lambda-handlers/v2/pkg/serviceerror"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func NewResponseWriter(w http.ResponseWriter, headers http.Header) *ResponseWriter {
	for k, v := range headers {
		if len(v) > 0 {
			w.Header().Set(k, v[0])
		}
	}

	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *ResponseWriter) Write(body []byte) (int, error) {
	w.body.Reset()
	bodyStr := string(body)
	if !helpers.IsOkRange(w.StatusCode()) && !helpers.IsValidJSONObject(bodyStr) {
		var decodedString string
		if err := json.Unmarshal([]byte(bodyStr), &decodedString); err == nil {
			bodyStr = decodedString
		}

		e := serviceerror.NewServiceError(
			serviceerror.GetServiceErrorCode(w.StatusCode()),
			serviceerror.GetServiceErrorCode(w.StatusCode()),
			bodyStr,
		)

		b, err := json.Marshal(e)
		if err != nil {
			slog.Error(err.Error())
			return 0, err
		}

		bodyStr = string(b)
	} else if !helpers.IsValidJSONObject(bodyStr) {
		var decodedString string
		err := json.Unmarshal([]byte(bodyStr), &decodedString)
		if err == nil {
			bodyStr = decodedString
		}
	}

	return w.body.Write([]byte(bodyStr))
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *ResponseWriter) Body() string {
	return w.body.String()
}

func (w *ResponseWriter) SetBody(b []byte) {
	w.body.Reset()
	w.body.Write(b)
}

func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *ResponseWriter) SetHttpResponseWriter() error {
	w.ResponseWriter.WriteHeader(w.statusCode)

	_, err := w.ResponseWriter.Write(w.body.Bytes())
	return err
}
