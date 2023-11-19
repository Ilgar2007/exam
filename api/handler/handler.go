package handler

import (
	"encoding/json"
	"exam/config"
	"exam/storage"
	"log"
	"net/http"
)

type Handler struct {
	cfg     *config.Config
	storage storage.StorageI
}

type Respons struct {
	Status      int         `json:"status"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}

func NewController(cfg *config.Config, strg storage.StorageI) *Handler {
	return &Handler{cfg: cfg, storage: strg}
}

func (h *Handler) handleResponse(w http.ResponseWriter, status int, data interface{}) {

	var description string
	switch code := status; {
	case code < 400:
		description = "success"
		sam, _ := json.MarshalIndent(Respons{
			Status:      status,
			Description: description,
			Data:        data,
		}, "", " ")
		log.Println("Status ", status, string(sam))

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(Respons{
			Status:      status,
			Description: description,
			Data:        data,
		})
	default:
		description = "error"
		log.Println(config.Error, "erro while:", Respons{
			Status:      status,
			Description: description,
			Data:        data,
		})

		if code == 500 {
			description = "Internal Server Error"
		}

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(Respons{
			Status:      status,
			Description: description,
			Data:        data,
		})
	}
}
