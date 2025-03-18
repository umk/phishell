package server

import (
	"encoding/json"
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

func parseResponseBody[T any](w http.ResponseWriter, res *http.Response) (T, bool) {
	var resp T
	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Error reading response body: "+err.Error(), http.StatusInternalServerError)
		return resp, false
	}
	if err := marshalx.UnmarshalJSONStruct(body, &resp); err != nil {
		http.Error(w, "Invalid response format: "+err.Error(), http.StatusInternalServerError)
		return resp, false
	}
	return resp, true
}

func propagateResponse(w http.ResponseWriter, res *http.Response, body any) {
	// Copy headers
	for k, v := range res.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	// Set status code
	w.WriteHeader(res.StatusCode)

	// Handle response body based on type.
	if body == nil {
		// do nothing
		return
	}
	if r, ok := body.(io.Reader); ok {
		io.Copy(w, r)
	} else {
		b, err := json.Marshal(body)
		if err != nil {
			http.Error(w, "Error marshaling response body: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}
