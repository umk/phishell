package server

import (
	"io"
	"net/http"

	"github.com/umk/phishell/util/marshalx"
)

func parseRequestBody[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var req T

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body: "+err.Error(), http.StatusBadRequest)
		return req, false
	}

	if err := marshalx.UnmarshalJSONStruct(body, &req); err != nil {
		http.Error(w, "Invalid request format: "+err.Error(), http.StatusBadRequest)
		return req, false
	}

	return req, true
}

// propagateResponse copies the response from an HTTP response to the response writer
func propagateResponse(w http.ResponseWriter, res *http.Response) {
	// Copy headers
	for k, v := range res.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	// Set status code
	w.WriteHeader(res.StatusCode)

	// Copy body
	defer res.Body.Close()
	io.Copy(w, res.Body)
}
