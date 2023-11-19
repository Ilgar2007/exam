package handler

import (
	"encoding/json"
	"exam/models"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cast"
)

func (h *Handler) Client(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.CreateClient(w, r)
	case "GET":
		var (
			id = r.URL.Query().Get("id")
		)
		if id != "" {
			h.ClientGetById(w, r)
		} else {
			h.ClientGetList(w, r)
		}
	case "DELETE":
		h.ClientDelete(w, r)
	case "PUT":
		h.ClientUpdate(w, r)
	default:
		fmt.Println("default")
	}
}

func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	var createClient models.CreateClient
	err := json.NewDecoder(r.Body).Decode(&createClient)

	if err != nil {
		h.handleResponse(w, 500, "error while unmarsh request: "+err.Error())
		return
	}

	resp, err := h.storage.Client().Create(&createClient)
	if err != nil {
		fmt.Println(createClient)
		h.handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.handleResponse(w, http.StatusCreated, resp)
}

func (h *Handler) ClientGetById(w http.ResponseWriter, r *http.Request) {
	var clId = r.URL.Query().Get("id")
	resp, err := h.storage.Client().GetByID(&models.ClientPrimaryKey{ID: clId})
	if err != nil {
		h.handleResponse(w, 500, "Client does not exist: "+err.Error())
		return
	}

	h.handleResponse(w, http.StatusOK, resp)
}

func (h *Handler) ClientGetList(w http.ResponseWriter, r *http.Request) {
	var offsetstr = r.URL.Query().Get("offset")
	var limitstr = r.URL.Query().Get("limit")
	var search = r.URL.Query().Get("search")

	var offset = cast.ToInt(offsetstr)
	var limit = cast.ToInt(limitstr)

	resp, err := h.storage.Client().GetList(&models.GetListClientRequest{Offset: offset, Limit: limit, Search: search})
	if err != nil {
		h.handleResponse(w, 500, "Client does not exist: "+err.Error())
		return
	}

	h.handleResponse(w, http.StatusOK, resp)
}
func (h *Handler) ClientDelete(w http.ResponseWriter, r *http.Request) {
	var clientId = r.URL.Query().Get("id")

	resp, err := h.storage.Client().Delete(&models.ClientPrimaryKey{ID: clientId})
	if err != nil {
		log.Println(err.Error())
		return
	}
	h.handleResponse(w, http.StatusOK, resp)
}

func (h *Handler) ClientUpdate(w http.ResponseWriter, r *http.Request) {
	var client models.UpdateClient

	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		log.Println(err.Error())
		return
	}
	resp, err := h.storage.Client().Update(&client)

	if err != nil {
		log.Println(err.Error())
		return
	}

	h.handleResponse(w, http.StatusOK, resp)
}
