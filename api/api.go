package api

import (
	"exam/api/handler"
	"exam/config"
	"exam/storage"
	"net/http"
)

func NewApi(cfg *config.Config, strg storage.StorageI) {
	handler := handler.NewController(cfg, strg)

	http.HandleFunc("/category", handler.Category)
	http.HandleFunc("/product", handler.Product)
	http.HandleFunc("/client", handler.Client)
	http.HandleFunc("/branch", handler.Branch)

}
