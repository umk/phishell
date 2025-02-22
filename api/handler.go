package api

import "net/http"

type handler struct {
}

func (*handler) CreateChatCompletion(w http.ResponseWriter, r *http.Request) {}

func (*handler) CreateCompletion(w http.ResponseWriter, r *http.Request) {}

func (*handler) CreateEmbedding(w http.ResponseWriter, r *http.Request) {}
