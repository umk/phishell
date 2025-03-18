package internal

import (
	"net/http"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	CreateChatCompletion(w http.ResponseWriter, r *http.Request)
	CreateEmbedding(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
	// HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// CreateChatCompletion operation middleware
func (siw *ServerInterfaceWrapper) CreateChatCompletion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateChatCompletion(w, r)
	}))

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateEmbedding operation middleware
func (siw *ServerInterfaceWrapper) CreateEmbedding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ctx = context.WithValue(ctx, ApiKeyAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateEmbedding(w, r)
	}))

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       *http.ServeMux
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:          si,
		ErrorHandlerFunc: options.ErrorHandlerFunc,
	}

	m.HandleFunc("POST "+options.BaseURL+"/chat/completions", wrapper.CreateChatCompletion)
	m.HandleFunc("POST "+options.BaseURL+"/embeddings", wrapper.CreateEmbedding)

	return m
}
